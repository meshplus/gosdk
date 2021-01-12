package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/meshplus/gosdk/common"
)

// EventType type
type EventType string

// SubscriptionID subscription id
type SubscriptionID string

const (
	// BLOCKEVENT block
	BLOCKEVENT EventType = "block"
	// SYSTEMSTATUSEVENT systemStatus
	SYSTEMSTATUSEVENT EventType = "systemStatus"
	// LOGSEVENT logs
	LOGSEVENT EventType = "logs"
)

const (
	// ReadBufferSize read buff size
	ReadBufferSize = 1024 * 256
	//WriteBufferSize write buff size
	WriteBufferSize = 1024 * 256
)

var globalWebSocketClient *WebSocketClient

// WsEventHandler web socket event handler
// note: if you unsubscribe a event, the OnClose() will never be called
// even if the connection is closed
type WsEventHandler interface {
	// when subscribe success
	OnSubscribe()
	// when unsubscribe
	OnUnSubscribe()
	// when receive notification
	OnMessage([]byte)
	// when connection closed
	OnClose()
}

type connectionWrapper struct {
	id            int
	conn          *websocket.Conn
	mutex         sync.Mutex // use to sync conn
	handler       WsEventHandler
	subIDCh       chan SubscriptionID
	eventHub      map[SubscriptionID]WsEventHandler
	subscriptions []Subscription
}

// WebSocketClient control the all APIs web socket related APIs
type WebSocketClient struct {
	conns   map[int]*connectionWrapper
	hrm     *httpRequestManager
	rwMutex sync.RWMutex // use to sync conns
}

// WebSocketNotification represents the notification data structure
type WebSocketNotification struct {
	Event        string          `json:"event"`
	Subscription SubscriptionID  `json:"subscription"`
	Data         json.RawMessage `json:"data"`
}

// Subscription -
type Subscription struct {
	Event          EventType      `json:"event"`
	SubscriptionID SubscriptionID `json:"subId"`
}

func (wscli *WebSocketClient) getConn(nodeIndex int) (*websocket.Conn, StdError) {
	nodeURL := wscli.hrm.nodes[nodeIndex].wsURL
	logger.Debug("web socket url:", nodeURL)

	header := make(http.Header)
	header["origin"] = []string{"haha"}
	if wscli.hrm.sendTcert {
		if wscli.hrm.tcm.cfca {
			header.Add("tcert", wscli.hrm.tcm.ecert)
			signature, err := wscli.hrm.tcm.sdkCert.Sign([]byte{})
			if err != nil {
				logger.Error("signature body error", err)
				return nil, NewSystemError(err)
			}
			header.Add("signature", common.Bytes2Hex(signature))
			header.Add("msg", common.Bytes2Hex([]byte{}))
		} else {
			header.Add("tcert", string(wscli.hrm.tcm.tcertPool[wscli.hrm.nodes[nodeIndex].url]))
			signature, err := wscli.hrm.tcm.uniqueCert.Sign([]byte{})
			if err != nil {
				logger.Error("signature body error,", err)
				return nil, NewSystemError(err)
			}
			header.Add("signature", common.Bytes2Hex(signature))
			header.Add("msg", common.Bytes2Hex([]byte{}))
		}
	}

	conn, resp, err := websocket.DefaultDialer.Dial(nodeURL, header)
	if err != nil || resp.StatusCode != http.StatusSwitchingProtocols {
		return nil, NewSystemError(errors.New("get webSocket connection error"))
	}

	return conn, nil
}

func (wscli *WebSocketClient) listen(wrapper *connectionWrapper) {
	// get a copy of connection, in case get nil pointer panic
	conn := wrapper.conn
	// heart beat
	go func(conn *websocket.Conn) {
		t := time.NewTicker(1 * time.Minute)
		defer t.Stop()
		for range t.C {
			err := conn.WriteControl(websocket.PingMessage, []byte("heart beat"), time.Now().Add(5*time.Second))
			if err != nil {
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
					logger.Errorf("web socket of node%d encountered an error: %v", wrapper.id, err)
				}
				if _, ok := err.(*websocket.CloseError); ok {
					return
				}
			}
		}
	}(conn)

	for {

		msgType, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				logger.Errorf("web socket of node%d encountered an error: %v", wrapper.id, err)
			}
			if _, ok := err.(*websocket.CloseError); ok {
				wscli.clearConn(wrapper.id)
				return
			}
		}

		jsonResponse := &JSONResponse{}

		err = json.Unmarshal(data, jsonResponse)
		if err != nil {
			logger.Errorf("web socket of node%d encountered an error: %v", wrapper.id, err)
			break
		}

		switch msgType {
		case websocket.TextMessage:
			logger.Debugf("[WEB SOCKET RESPONSE]: %s", string(data))

			// subscription result
			if jsonResponse.Message != "" {
				if jsonResponse.Message == "SUCCESS" {
					// if not "true" or "false"
					if len(jsonResponse.Result) != 4 && len(jsonResponse.Result) != 5 {
						var subID SubscriptionID

						err := json.Unmarshal(jsonResponse.Result, &subID)
						if err != nil {
							logger.Errorf("web socket of node%d encountered an error: %v", wrapper.id, err)
							break
						}

						go func() {
							// inform user of the subID
							wrapper.subIDCh <- SubscriptionID(subID)
							wrapper.eventHub[SubscriptionID(subID)] = wrapper.handler
							// notify user
							wrapper.eventHub[SubscriptionID(subID)].OnSubscribe()
						}()
					}
				}
			} else {
				// event notification
				var notification WebSocketNotification
				err := json.Unmarshal(jsonResponse.Result, &notification)
				if err != nil {
					logger.Errorf("web socket of node %d encountered an error: %v", wrapper.id, err)
					break
				}
				// if the callback has been removed, try to unsubscribe again
				if _, ok := wrapper.eventHub[notification.Subscription]; !ok {
					//nolint
					go wscli.UnSubscribe(notification.Subscription)
				} else {
					// notify user
					go wrapper.eventHub[notification.Subscription].OnMessage(notification.Data)
				}
			}
		default:
			break
		}
	}

}

func (wscli *WebSocketClient) checkIndex(index int) bool {
	return 0 <= index && index <= len(wscli.hrm.nodes)
}

func (wscli *WebSocketClient) clearConn(nodeIndex int) {
	if wscli.conns[nodeIndex] != nil {
		// clear subIDs
		wscli.conns[nodeIndex].subscriptions = make([]Subscription, 0)
		// close channel
		close(wscli.conns[nodeIndex].subIDCh)
		// clear the connection
		wscli.conns[nodeIndex] = nil
	}
}

func (wscli *WebSocketClient) getCloseHandler(wrapper *connectionWrapper) func(code int, text string) error {
	return func(code int, text string) error {
		// notify user
		for _, h := range wrapper.eventHub {
			h.OnClose()
		}

		// clear the connection and callback
		wscli.clearConn(wrapper.id)

		return nil
	}
}

// SubscribeForProposal is used to subscribe logs about proposal of the specific node, the user-defined callback
// will be called when proposal status been changed.
// note: nodeIndex start from 1
func (wscli *WebSocketClient) SubscribeForProposal(nodeIndex int, eventHandler WsEventHandler) (SubscriptionID, StdError) {
	method := "bvm_subscribe"
	params := make([]interface{}, 1)
	params[0] = "configProposalSubscribe"
	return wscli.subscribe(nodeIndex, eventHandler, method, params, LOGSEVENT)
}

// Subscribe is used to subscribe event(s) of the specific node, the user-defined callback
// will be called when events that fulfill the filters occurred.
// note: nodeIndex start from 1
func (wscli *WebSocketClient) Subscribe(nodeIndex int, filter EventFilter, eventHandler WsEventHandler) (SubscriptionID, StdError) {
	method := "sub_subscribe"
	params := make([]interface{}, 2)
	params[0] = filter.GetEventType()
	params[1] = filter.Serialize()
	return wscli.subscribe(nodeIndex, eventHandler, method, params, filter.GetEventType())
}

func (wscli *WebSocketClient) subscribe(nodeIndex int, eventHandler WsEventHandler, method string, params []interface{}, eventType EventType) (SubscriptionID, StdError) {
	if !wscli.checkIndex(nodeIndex) {
		return "", NewSystemError(fmt.Errorf("node index out of range, suppose to be in [0, %d)", len(wscli.hrm.nodes)))
	}

	var (
		stdErr StdError
		conn   *websocket.Conn
	)
	nodeIndex = nodeIndex - 1

	// lock conns
	wscli.rwMutex.Lock()
	defer wscli.rwMutex.Unlock()

	// lazy init
	if wscli.conns[nodeIndex] == nil {
		wscli.conns[nodeIndex] = &connectionWrapper{
			mutex:         sync.Mutex{},
			subIDCh:       make(chan SubscriptionID),
			id:            nodeIndex,
			subscriptions: make([]Subscription, 0),
			eventHub:      make(map[SubscriptionID]WsEventHandler),
		}
		if wscli.conns[nodeIndex].conn, stdErr = wscli.getConn(nodeIndex); stdErr != nil {
			return "", stdErr
		}
		// init when connect success
		wscli.conns[nodeIndex].conn.SetCloseHandler(wscli.getCloseHandler(wscli.conns[nodeIndex]))
		wscli.conns[nodeIndex].conn.SetPingHandler(wscli.pingHandler)
		wscli.conns[nodeIndex].conn.SetPongHandler(wscli.pongHandler)

		go wscli.listen(wscli.conns[nodeIndex])
	}

	conn = wscli.conns[nodeIndex].conn

	jsonReq := JSONRequest{
		Method:    method,
		Version:   JSONRPCVersion,
		ID:        1,
		Namespace: wscli.hrm.namespace,
		Params:    params,
	}

	req, err := json.Marshal(jsonReq)
	if err != nil {
		return "", NewSystemError(err)
	}

	logger.Debugf("[WEB SOCKET REQUEST]: %s", string(req))

	// before request
	wscli.conns[nodeIndex].handler = eventHandler

	// lock conn
	wscli.conns[nodeIndex].mutex.Lock()
	err = conn.WriteMessage(websocket.TextMessage, req)
	if err != nil {
		wscli.conns[nodeIndex].mutex.Unlock()
		return "", NewSystemError(err)
	}
	wscli.conns[nodeIndex].mutex.Unlock()

	subID, ok := <-wscli.conns[nodeIndex].subIDCh
	if !ok {
		return "", NewSystemError(fmt.Errorf("the connection of the node%d is closed, please try again", nodeIndex))
	}

	// store eventType=>subID
	wscli.conns[nodeIndex].subscriptions =
		append(wscli.conns[nodeIndex].subscriptions, Subscription{Event: eventType, SubscriptionID: subID})

	return subID, nil
}

// CloseConn is used to close the connection of a specific node
// note: nodeIndex start from 1
func (wscli *WebSocketClient) CloseConn(nodeIndex int) StdError {
	if !wscli.checkIndex(nodeIndex) {
		return NewSystemError(fmt.Errorf("node index out of range, suppose to be in [0, %d)", len(wscli.hrm.nodes)))
	}

	nodeIndex = nodeIndex - 1

	wscli.rwMutex.Lock()
	defer wscli.rwMutex.Unlock()

	// try to close connection
	if wscli.conns[nodeIndex] != nil {

		wscli.conns[nodeIndex].mutex.Lock()
		err := wscli.conns[nodeIndex].conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "ok"))
		if err != nil {
			wscli.conns[nodeIndex].mutex.Unlock()
			return NewSystemError(err)
		}
		wscli.conns[nodeIndex].mutex.Unlock()

		wscli.clearConn(nodeIndex)
	}

	return nil
}

// UnSubscribe is used to unsubscribe a event by subID and user will not be notified by
// the event once the method called
func (wscli *WebSocketClient) UnSubscribe(id SubscriptionID) StdError {
	jsonReq := &JSONRequest{
		ID:        1,
		Namespace: wscli.hrm.namespace,
		Version:   JSONRPCVersion,
		Method:    "sub_unsubscribe",
		Params:    []interface{}{id},
	}

	req, err := json.Marshal(jsonReq)
	if err != nil {
		return NewSystemError(err)
	}

	for k := range wscli.conns {
		for k1 := range wscli.conns[k].eventHub {
			if k1 == id {
				logger.Debugf("[WEB SOCKET REQUEST]: %s", string(req))

				wscli.rwMutex.RLock()
				defer wscli.rwMutex.RUnlock()

				// try to unsubscribe
				if wscli.conns[k] != nil {
					wscli.conns[k].mutex.Lock()
					if err := wscli.conns[k].conn.WriteMessage(websocket.TextMessage, req); err != nil {
						wscli.conns[k].mutex.Unlock()
						return NewSystemError(err)
					}
					wscli.conns[k].mutex.Unlock()

					// notify user
					go wscli.conns[k].eventHub[id].OnUnSubscribe()
					// clear callback
					delete(wscli.conns[k].eventHub, id)
					// clear subID
					for i := range wscli.conns[k].subscriptions {
						if wscli.conns[k].subscriptions[i].SubscriptionID == id {
							wscli.conns[k].subscriptions = append(wscli.conns[k].subscriptions[:i], wscli.conns[k].subscriptions[i+1:]...)
							break
						}
					}
				}
				return nil
			}
		}
	}

	return nil
}

// GetAllSubscription get all subscriptions of a specific node
// note: nodeIndex start from 1
func (wscli *WebSocketClient) GetAllSubscription(nodeIndex int) ([]Subscription, StdError) {
	if !wscli.checkIndex(nodeIndex) {
		return nil, NewSystemError(fmt.Errorf("node index out of range, suppose to be in [0, %d)", len(wscli.hrm.nodes)))
	}

	nodeIndex = nodeIndex - 1

	wscli.rwMutex.RLock()
	defer wscli.rwMutex.RUnlock()

	if wscli.conns[nodeIndex] == nil {
		return make([]Subscription, 0), nil
	}

	return wscli.conns[nodeIndex].subscriptions, nil
}

func (wscli *WebSocketClient) pingHandler(appData string) error {
	logger.Debugf("[WEB SOCKET PING]: %s", appData)
	return nil
}

func (wscli *WebSocketClient) pongHandler(appData string) error {
	logger.Debugf("[WEB SOCKET PONG]: %s", appData)
	return nil
}

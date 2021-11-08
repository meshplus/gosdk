package rpc

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/meshplus/gosdk/common"
	"github.com/terasum/viper"
)

type RequestType string

const (
	GENERAL  RequestType = "GENERAL"
	UPLOAD   RequestType = "UPLOAD"
	DOWNLOAD RequestType = "DOWNLOAD"
)

// Node is used to contain node info
type Node struct {
	url    string
	wsURL  string
	status bool
}

func newNode(url string, rpcPort string, wsPort string, isHTTPS bool) (node *Node) {
	var scheme string

	if isHTTPS {
		scheme = "https://"
	} else {
		scheme = "http://"
	}

	node = &Node{
		url:    scheme + url + ":" + rpcPort,
		wsURL:  "ws://" + url + ":" + wsPort,
		status: true,
	}
	return node
}

// NewNode create a new node
func NewNode(url string, rpcPort string, wsPort string) (node *Node) {
	return newNode(url, rpcPort, wsPort, false)
}

// httpRequestManager is used to manager node and http request
type httpRequestManager struct {
	nodes      []*Node
	nodeIndex  int
	client     *http.Client
	namespace  string
	sendTcert  bool
	isHTTP     bool
	tcm        *TCertManager
	reConnTime int64
	txVersion  string
}

// newHTTPRequestManager is used to construct httpRequestManager
func newHTTPRequestManager(vip *viper.Viper, confRootPath string, txVersion string) (hrm *httpRequestManager) {
	var (
		namespace string
		urls      []string
		rpcPorts  []string
		wsPorts   []string
		isHTTPS   bool
		tcm       *TCertManager
	)

	namespace = vip.GetString(common.NamespaceConf)

	urls = vip.GetStringSlice(common.JSONRPCNodes)
	logger.Debugf("[CONFIG]: %s = %v", common.JSONRPCNodes, urls)

	rpcPorts = vip.GetStringSlice(common.JSONRPCPorts)
	logger.Debugf("[CONFIG]: %s = %v", common.JSONRPCPorts, rpcPorts)

	wsPorts = vip.GetStringSlice(common.WebSocketPorts)
	logger.Debugf("[CONFIG]: %s = %v", common.WebSocketPorts, wsPorts)

	isHTTPS = vip.GetBool(common.SecurityHttps)
	logger.Debugf("[CONFIG]: %s = %v", common.SecurityHttps, isHTTPS)

	reConnTime := vip.GetInt64(common.ReConnectTime)

	var nodes = make([]*Node, len(urls))

	for i, url := range urls {
		nodes[i] = newNode(url, rpcPorts[i], wsPorts[i], isHTTPS)
	}

	sendTcert := vip.GetBool(common.PrivacySendTcert)
	logger.Debugf("[CONFIG]: sendTcert = %v", sendTcert)

	tcm = NewTCertManager(vip, confRootPath)

	httpRequestManager := &httpRequestManager{
		nodes:      nodes,
		nodeIndex:  0,
		client:     newHTTPClient(vip, confRootPath),
		namespace:  namespace,
		sendTcert:  sendTcert,
		tcm:        tcm,
		isHTTP:     isHTTPS,
		reConnTime: reConnTime,
		txVersion:  txVersion,
	}

	if sendTcert && !vip.GetBool(common.PrivacyCfca) && !isFlato(txVersion) {
		tcm.tcertPool = make(map[string]TCert)
		for _, node := range nodes {
			tcert, err := httpRequestManager.getTCert(node.url)
			if err != nil {
				// if getTCert's method is not exist, means platform is flato
				if err.Code() == MethodNotExistOrInvalidCode {
					return
				}
				logger.Error("can not get tcert from ", node.url, err)
				return
			}
			tcm.tcertPool[node.url] = TCert(tcert)
		}
	}

	return httpRequestManager
}

func newHTTPClient(vip *viper.Viper, confRootPath string) *http.Client {
	logger.Debugf("[CONFIG]: https = %v", vip.GetBool(common.SecurityHttps))

	maxIdleConns := vip.GetInt(common.MaxIdleConns)
	maxIdleConnsPerHost := vip.GetInt(common.MaxIdleConnsPerHost)

	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		DisableKeepAlives:     false,
		MaxIdleConns:          maxIdleConns,
		MaxIdleConnsPerHost:   maxIdleConnsPerHost,
		TLSHandshakeTimeout:   10 * time.Second, //TLS安全连接握手超时时间
		ExpectContinueTimeout: 1 * time.Second,  //发送完请求到接收到响应头的超时时间
	}

	if vip.GetBool(common.SecurityHttps) {
		pool := x509.NewCertPool()

		tlscaPath := strings.Join([]string{confRootPath, vip.GetString(common.SecurityTlsca)}, "/")
		tlspeerCertPath := strings.Join([]string{confRootPath, vip.GetString(common.SecurityTlspeerCert)}, "/")
		tlspeerCertPrivPath := strings.Join([]string{confRootPath, vip.GetString(common.SecurityTlspeerPriv)}, "/")

		caCrt, err := ioutil.ReadFile(tlscaPath)
		if err != nil {
			panic(fmt.Sprintf("read tlsCA from %s failed", tlscaPath))
		}

		pool.AppendCertsFromPEM(caCrt)

		cliCrt, err := tls.LoadX509KeyPair(tlspeerCertPath, tlspeerCertPrivPath)
		if err != nil {
			panic(fmt.Sprintf("read tlspeerCert from %s and %s failed", tlspeerCertPath, tlspeerCertPrivPath))
		}

		tr.TLSClientConfig = &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
			ServerName:   "hyperchain.cn",
		}
	}

	return &http.Client{Transport: tr}
}

func defaultHTTPRequestManager() *httpRequestManager {
	return &httpRequestManager{
		namespace: DefaultNamespace,
		nodes:     make([]*Node, 0),
		nodeIndex: 0,
		client:    &http.Client{},
		sendTcert: false,
		tcm:       nil,
		isHTTP:    false,
	}
}

func post(url string, body []byte) (*http.Request, StdError) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	return req, NewGetResponseError(err)
}

func addHeaders(req *http.Request, extraHeaders map[string]string) {
	if extraHeaders != nil {
		for k, v := range extraHeaders {
			req.Header.Add(k, v)
		}
	}
}

// SyncRequest function is used to send http request
func (hrm *httpRequestManager) SyncRequest(body []byte) ([]byte, StdError) {
	randomURL, stdErr := hrm.randomURL()
	if stdErr != nil {
		return nil, stdErr
	}

	return hrm.SyncRequestSpecificURL(body, randomURL, GENERAL, nil, nil)
}

// SyncRequestSpecificURL is used to post request to specific url
func (hrm *httpRequestManager) SyncRequestSpecificURL(body []byte, url string, requestType RequestType, extraHeaders map[string]string, rwSeeker io.ReadWriteSeeker) ([]byte, StdError) {
	var req *http.Request
	var stdErr StdError
	switch requestType {
	case DOWNLOAD:
		req, stdErr = post(url, body)
		if stdErr != nil {
			return nil, stdErr
		}
		addHeaders(req, extraHeaders)
		req.Header.Add("content-type", "application/octet-stream")
	case UPLOAD:
		var err error
		req, err = http.NewRequest("POST", url, rwSeeker)
		if err != nil {
			return nil, NewSystemError(err)
		}
		addHeaders(req, extraHeaders)
		req.Header.Add("content-type", "application/octet-stream")
	case GENERAL:
		fallthrough
	default:
		req, stdErr = post(url, body)
		if stdErr != nil {
			return nil, stdErr
		}
	}

	if hrm.sendTcert {
		if isFlato(hrm.txVersion) || hrm.tcm.cfca {
			signature, sysErr := hrm.tcm.sdkCert.Sign(body)
			if sysErr != nil {
				logger.Error("sign error", sysErr)
				return nil, NewSystemError(sysErr)
			}
			req.Header.Add("tcert", hrm.tcm.ecert)
			req.Header.Add("signature", common.Bytes2Hex(signature))
		} else {
			signature, err := hrm.tcm.uniqueCert.Sign(body)
			if err != nil {
				logger.Error("signature body error,", err)
				return nil, NewSystemError(err)
			}
			req.Header.Add("tcert", string(hrm.tcm.tcertPool[url]))
			req.Header.Add("signature", common.Bytes2Hex(signature))
		}
	}

	logger.Debug("[URL]:", url)
	logger.Debug("[REQUEST]:", string(body))

	resp, sysErr := hrm.client.Do(req)
	if sysErr != nil {
		return nil, NewGetResponseError(sysErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// 200
		if requestType == DOWNLOAD && resp.Header.Get("Content-Type") == "application/octet-stream" {
			strPos := extraHeaders["pos"]
			var pos int64
			if strPos != "" {
				var err error
				pos, err = strconv.ParseInt(strPos, 10, 64)
				if err != nil {
					logger.Warning("pos convert failed, use default value 0")
				}
			}
			fsErr := streamFileStorage(rwSeeker, resp.Body, pos)
			if fsErr != nil {
				return nil, NewSystemError(fsErr)
			}
			return newFakeJSONResponse(0, "download success", hrm.txVersion), nil
		} else {
			ret, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, NewSystemError(err)
			}
			logger.Debug("[RESPONSE]:", string(ret))
			return ret, nil
		}
	} else if !isTemporaryError(resp.StatusCode) {
		return nil, NewHttpResponseError(resp.StatusCode, resp.Status)
	}

	// 请求异常返回，重连节点
	hrm.ReConnectNode(hrm.nodeIndex)

	return nil, NewGetResponseError(errors.New("http failed " + resp.Status))
}

func isTemporaryError(code int) bool {
	return code >= 500 && // fast short cut
		code != http.StatusNotImplemented &&
		code != http.StatusVariantAlsoNegotiates &&
		code != http.StatusHTTPVersionNotSupported &&
		code != http.StatusNotExtended
}

func (hrm *httpRequestManager) getTCert(url string) (string, StdError) {
	rawReq := &JSONRequest{
		Method:    "cert_getTCert",
		Version:   JSONRPCVersion,
		ID:        1,
		Namespace: hrm.namespace,
	}
	uniqPub, sysErr := ioutil.ReadFile(hrm.tcm.uniquePubPath)
	if sysErr != nil {
		return "", NewSystemError(sysErr)
	}
	param := newMapParam("pubkey", common.Bytes2Hex(uniqPub)).Serialize()
	rawReq.Params = []interface{}{param}

	body, sysErr := json.Marshal(rawReq)
	if sysErr != nil {
		return "", NewSystemError(sysErr)
	}

	req, stdErr := post(url, body)
	if stdErr != nil {
		return "", stdErr
	}

	signature, sysErr := hrm.tcm.sdkCert.Sign(body)
	if sysErr != nil {
		return "", NewSystemError(sysErr)
	}
	req.Header.Add("tcert", hrm.tcm.ecert)
	req.Header.Add("signature", common.Bytes2Hex(signature))
	req.Header.Add("msg", common.Bytes2Hex(body))

	logger.Debug("[URL]:", url)
	logger.Debug("[REQUEST]:", string(body))

	resp, sysErr := hrm.client.Do(req)
	if sysErr != nil {
		return "", NewGetResponseError(sysErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		ret, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", NewSystemError(err)
		}
		logger.Debug("[RESPONSE]:", string(ret))

		var resp *JSONResponse
		if sysErr = json.Unmarshal(ret, &resp); sysErr != nil {
			return "", NewSystemError(sysErr)
		}

		if resp.Code != SuccessCode {
			return "", NewServerError(resp.Code, resp.Message)
		}

		var tcert TCertResponse
		if err := json.Unmarshal(resp.Result, &tcert); err != nil {
			return "", NewSystemError(err)
		}
		return tcert.TCert, nil
	}

	return "", NewGetResponseError(errors.New("http failed " + resp.Status))
}

func (hrm *httpRequestManager) randomURL() (url string, err StdError) {
	nodeNum := len(hrm.nodes)
	randomNum := nodeNum * 2
	for randomNum > 0 {
		hrm.nodeIndex = common.RandInt(nodeNum)
		if hrm.nodes[hrm.nodeIndex].status {
			return hrm.nodes[hrm.nodeIndex].url, nil
		}
		randomNum--
	}
	logger.Error("All nodes are bad, please check it! Now retry to connect all nodes.")

	//if random fail, try round
	for i := 0; i < nodeNum; i++ {
		hrm.nodeIndex = (hrm.nodeIndex + 1) % nodeNum
		if hrm.nodes[hrm.nodeIndex].status {
			return hrm.nodes[hrm.nodeIndex].url, nil
		}
	}

	return "", NewGetResponseError(errors.New("all nodes are bad, please check it"))
}

// getNodeURL get the url of the node
func (hrm *httpRequestManager) getNodeURL(nodeID int) (url string, err StdError) {
	if nodeID < 0 || nodeID > len(hrm.nodes) {
		return "", NewGetResponseError(errors.New("node id is out of nodes size"))
	}
	if nodeID == 0 {
		return hrm.randomURL()
	}
	if !hrm.nodes[nodeID-1].status {
		return "", NewGetResponseError(errors.New(fmt.Sprintf("node %d is bad, please check it", nodeID)))
	}
	return hrm.nodes[nodeID-1].url, nil
}

// ReConnectNode is used to reconnect the node by index
func (hrm *httpRequestManager) ReConnectNode(nodeIndex int) {
	hrm.nodes[nodeIndex].status = false
	url := hrm.nodes[nodeIndex].url
	req := &JSONRequest{
		Method:    "node_getNodes",
		Version:   JSONRPCVersion,
		ID:        1,
		Namespace: hrm.namespace,
	}
	body, err := json.Marshal(req)
	if err != nil {
		logger.Error(NewSystemError(err).String())
	}

	go func() {
		request, err := post(url, body)
		if err != nil {
			logger.Error(err.String())
		}

		for {
			response, err := hrm.client.Do(request)
			if err != nil {
				logger.Error(NewSystemError(err).String())
			}

			if response != nil && response.StatusCode == http.StatusOK {
				b, _ := ioutil.ReadAll(response.Body)
				logger.Debug("reconnection node body: ", string(b))
				response.Body.Close()
				hrm.nodes[nodeIndex].status = true
				logger.Info("node " + hrm.nodes[nodeIndex].url + " Reconnect Success!")
				return
			}
			response.Body.Close()
			logger.Info("node " + hrm.nodes[nodeIndex].url + " Reconnect failed, will try one second later")
			time.Sleep(time.Millisecond * time.Duration(hrm.reConnTime))
		}
	}()

}

func isFlato(TxVersion string) bool {
	return TxVersion != DefaultTxVersion
}

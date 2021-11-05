// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package abi

import (
	"bytes"
	"encoding/json"
	"strings"

	"fmt"
	"io"
)

// The ABI holds information about a contract's context and available
// invokable methods. It will allow you to type check function calls and
// packs data accordingly.
type ABI struct {
	Constructor Method
	Methods     map[string]Method
	Events      map[string]Event

	// Fields that record overloaded methods and events
	allMethods map[string]Method
	allEvents  map[string]Event
}

// JSON returns a parsed ABI interface and error if it failed.
func JSON(reader io.Reader) (ABI, error) {
	dec := json.NewDecoder(reader)

	var abi ABI
	if err := dec.Decode(&abi); err != nil {
		return ABI{}, err
	}

	return abi, nil
}

// GetMethod is used to get method using method signature
func (abi ABI) GetMethod(methodSig string) (*Method, error) {
	if !strings.HasSuffix(methodSig, ")") {
		// To get method only with its name, it will return the method with no argument by default
		if method, ok1 := abi.allMethods[methodSig+"()"]; ok1 {
			return &method, nil
		} else {
			if lastMatchedMethod, ok2 := abi.Methods[methodSig]; ok2 { // If the method with no argument doesn't exist, return the last matched method
				return &lastMatchedMethod, nil
			} else {
				return nil, fmt.Errorf("no method with signature %s found", methodSig)
			}
		}
	} else {
		if method, ok := abi.allMethods[methodSig]; ok {
			return &method, nil
		} else {
			return nil, fmt.Errorf("no method with signature %s found", methodSig)
		}
	}
}

// GetEvent is used to get event using event signature
func (abi ABI) GetEvent(eventSig string) (*Event, error) {
	if !strings.HasSuffix(eventSig, ")") {
		// To get event only with its name, it will return the event with no argument by default
		if event, ok1 := abi.allEvents[eventSig+"()"]; ok1 {
			return &event, nil
		} else {
			if lastMatchedEvent, ok2 := abi.Events[eventSig]; ok2 { // If the event with no argument doesn't exist, return the last matched event
				return &lastMatchedEvent, nil
			} else {
				return nil, fmt.Errorf("no event with signature %s found", eventSig)
			}
		}
	} else {
		if event, ok := abi.allEvents[eventSig]; ok {
			return &event, nil
		} else {
			return nil, fmt.Errorf("no event with signature %s found", eventSig)
		}
	}
}

// GetPayload is used to convert func and args to binary
func (abi ABI) GetPayload(funcName string, args ...string) string {
	var payload string
	var typeArgs []string

	method, err := abi.GetMethod(funcName)
	if err != nil {
		logger.Error("no method found: ", err)
		return ""
	}
	for _, argument := range method.Inputs {
		typeArgs = append(typeArgs, argument.Type.String())
	}

	funcNameABI := FuncSelector(funcName, typeArgs)
	for i := 0; i < len(args); i++ {
		solcBin := TypeConversion(args[i], typeArgs[i])
		payload += solcBin
	}
	payload = strings.TrimPrefix(funcNameABI+payload, "0x")

	return payload
}

// Pack the given method name to conform the ABI. Method call's data
// will consist of method_id, args0, arg1, ... argN. Method id consists
// of 4 bytes and arguments are all 32 bytes.
// Method ids are created from the first 4 bytes of the hash of the
// methods string signature. (signature = baz(uint32,string32))
func (abi ABI) Pack(name string, args ...interface{}) ([]byte, error) {
	// Fetch the ABI of the requested method
	if name == "" {
		// constructor
		arguments, err := abi.Constructor.Inputs.Pack(args...)
		if err != nil {
			return nil, err
		}
		return arguments, nil

	}
	method, err := abi.GetMethod(name)
	if err != nil {
		return nil, err
	}

	arguments, err := method.Inputs.Pack(args...)
	if err != nil {
		return nil, err
	}
	// Pack up the method ID too if not a constructor and return
	return append(method.Id(), arguments...), nil
}

// Unpack output in v according to the abi specification
func (abi ABI) Unpack(v interface{}, name string, output []byte) (err error) {
	if len(output) == 0 {
		return fmt.Errorf("abi: unmarshalling empty output")
	}
	// since there can't be naming collisions with contracts and events,
	// we need to decide whether we're calling a method or an event
	if method, err := abi.GetMethod(name); err == nil {
		if len(output)%32 != 0 {
			return fmt.Errorf("abi: improperly formatted output")
		}
		return method.Outputs.Unpack(v, output)
	} else if event, err := abi.GetEvent(name); err == nil {
		return event.Inputs.Unpack(v, output)
	}
	return fmt.Errorf("abi: could not locate named method or event")
}

// UnmarshalJSON implements json.Unmarshaler interface
func (abi *ABI) UnmarshalJSON(data []byte) error {
	var fields []struct {
		Type      string
		Name      string
		Constant  bool
		Anonymous bool
		Inputs    []Argument
		Outputs   []Argument
	}

	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	abi.Methods = make(map[string]Method)
	abi.Events = make(map[string]Event)
	abi.allMethods = make(map[string]Method)
	abi.allEvents = make(map[string]Event)
	for _, field := range fields {
		switch field.Type {
		case "constructor":
			abi.Constructor = Method{
				Inputs: field.Inputs,
			}
			// empty defaults to function according to the abi spec
		case "function", "":
			method := Method{
				Name:    field.Name,
				Const:   field.Constant,
				Inputs:  field.Inputs,
				Outputs: field.Outputs,
			}
			methodSig := method.Sig()
			abi.Methods[field.Name] = method
			abi.allMethods[methodSig] = method
		case "event":
			event := Event{
				Name:      field.Name,
				Anonymous: field.Anonymous,
				Inputs:    field.Inputs,
			}
			eventSig := event.Sig()
			abi.Events[field.Name] = event
			abi.allEvents[eventSig] = event
		}
	}

	return nil
}

// MethodById looks up a method by the 4-byte id
// returns nil if none found
func (abi *ABI) MethodById(sigdata []byte) (*Method, error) {
	if len(sigdata) < 4 {
		return nil, fmt.Errorf("data too short (% bytes) for abi method lookup", len(sigdata))
	}
	for _, method := range abi.Methods {
		if bytes.Equal(method.Id(), sigdata[:4]) {
			return &method, nil
		}
	}
	return nil, fmt.Errorf("no method with id: %#x", sigdata[:4])
}

package command

import (
	"encoding/binary"
	"fmt"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// IAction if Action interface
type IAction interface {
	Execute() (interface{}, error)
	XXX_MapID() uint16
	XXX_MapName() string
}

// IResponse if Response interface
type IResponse interface {
	XXX_MapID() uint16
	XXX_MapName() string
}

// IMap map id and object
type IMap interface {
	IDToObject(id uint16) interface{}
	IDCount() uint16
}

// Dispatch manage action,handler,response
type Dispatch struct {
	Map IMap
}

// ErrCommandParsing fire when decode command has error, it 's client 's fault
var ErrCommandParsing = errors.New("command parsing error")

// Route get action from httpRequest and write response to httpResponse
// write error text if some thing is wrong
func (dp *Dispatch) Route(bytes []byte) ([]byte, error) {
	//bytes is command contain [proto,id], id is 2 bytes
	_, action, err := dp.decodeCommand(bytes)
	if err == ErrCommandParsing {
		return nil, err
	}
	if err != nil {
		ErrCommandParsing = errors.Wrap(err, "command parsing error")
		return nil, ErrCommandParsing
	}

	responseID, response, err := dp.handle(action)
	if err != nil {
		return nil, errors.Wrap(err, "failed to handle action")
	}

	startTime := time.Now()
	fmt.Printf("execute %v(%v bytes)", action.(IAction).XXX_MapName(), len(bytes))
	var returnBytes []byte
	if response != nil {
		returnBytes, err = dp.encodeCommand(responseID, response)
		fmt.Printf("return %v(%v bytes)", response.(IResponse).XXX_MapName(), len(returnBytes))
	}
	duration := time.Now().Sub(startTime)
	ms := duration.Nanoseconds() / 10000000
	fmt.Printf(", %v ms\n", ms)
	return returnBytes, err
}

//fastAppend provide better performance than append
func (dp *Dispatch) fastAppend(bytes1 []byte, bytes2 []byte) []byte {
	//return append(bytes1[:], bytes2[:]...)
	totalLen := len(bytes1) + len(bytes2)
	tmp := make([]byte, totalLen)
	i := copy(tmp, bytes1)
	copy(tmp[i:], bytes2)
	return tmp
}

//protoFromBuffer read proto message from buffer
func (dp *Dispatch) protoFromBuffer(id uint16, bytes []byte) (interface{}, error) {
	obj := dp.Map.IDToObject(id)
	if obj == nil {
		return nil, errors.New(fmt.Sprintf("failed to map id %v", id))
	}
	err := proto.Unmarshal(bytes, obj.(proto.Message))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode protobuf")
	}
	return obj, nil
}

//protoToBuffer write proto message to buffer
func (dp *Dispatch) protoToBuffer(obj interface{}) ([]byte, error) {
	if obj == nil {
		return nil, errors.New("obj nil")
	}

	bytes, err := proto.Marshal(obj.(proto.Message))
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode protobuf")
	}
	return bytes, nil
}

// handle send action to handler and get response
func (dp *Dispatch) handle(action interface{}) (uint16, interface{}, error) {
	responseInterface, err := action.(IAction).Execute()
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to execute action")
	}
	if responseInterface == nil {
		return 0, nil, nil
	}
	response := responseInterface.(IResponse)
	return response.XXX_MapID(), response, nil
}

// encodeCommand, comand is array contain [protobuf,id]
func (dp *Dispatch) encodeCommand(id uint16, proto interface{}) ([]byte, error) {
	bytes, err := dp.protoToBuffer(proto)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert proto to buffer")
	}
	idBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(idBytes, id)
	return dp.fastAppend(bytes, idBytes), nil
}

// decodeCommand, comand is array contain [protobuf,id]
func (dp *Dispatch) decodeCommand(bytes []byte) (uint16, interface{}, error) {
	bytesLen := len(bytes)
	protoBytes := bytes[:bytesLen-2]
	idBytes := bytes[bytesLen-2:]
	id := binary.LittleEndian.Uint16(idBytes)
	if id > dp.Map.IDCount() {
		ErrCommandParsing = errors.New(fmt.Sprintf("id %v out of range, the max id allowed is %v", id, dp.Map.IDCount()))
		return 0, nil, ErrCommandParsing
	}

	protoInterface, err := dp.protoFromBuffer(id, protoBytes)
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to convert buffer to proto")
	}
	return id, protoInterface, nil
}

package command

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/piyuo/go-libsrv/app"

	proto "github.com/golang/protobuf/proto"
	shared "github.com/piyuo/go-libsrv/command/shared"
	log "github.com/piyuo/go-libsrv/log"
	tools "github.com/piyuo/go-libsrv/tools"

	"github.com/pkg/errors"
)

// Action interface
type Action interface {
	Main(ctx context.Context) (interface{}, error)
	//  1 to 100 is shared command id between all service, 101 to 65,535 is valid service id
	XXX_MapID() uint16
	XXX_MapName() string
}

// Response interface
type Response interface {
	//  1 to 32,767 is valid service id,-1 to -32,768 is shared id between all service
	XXX_MapID() uint16
	XXX_MapName() string
}

// IMap map id and object
type IMap interface {
	NewObjectByID(id uint16) interface{}
}

// Dispatch manage action,handler,response
type Dispatch struct {
	Map IMap
}

// Route get action from httpRequest and write response to httpResponse
//
// write http error text if some thing went wrong
func (dp *Dispatch) Route(ctx context.Context, bytes []byte) ([]byte, error) {
	//bytes is command contain [proto,id], id is 2 bytes
	_, action, err := dp.decodeCommand(bytes)
	if err != nil {
		return nil, err
	}
	commandLog := fmt.Sprintf("exec %v (%v bytes), ", action.(Action).XXX_MapName(), len(bytes))
	responseID, response, err := dp.timeExecuteAction(ctx, action)
	if err != nil {
		return nil, err
	}
	var returnBytes []byte
	returnBytes, err = dp.encodeCommand(responseID, response)
	if err != nil {
		//commandLog += fmt.Sprintf("failed with %v , %v ms\n", err.Error(), ms)
		commandLog += fmt.Sprintf("failed with %v\n", err.Error())
		log.Debug(ctx, here, commandLog)
		return nil, err
	}
	commandLog += fmt.Sprintf("return %v (%v bytes)\n", betterResponseName(responseID, response), len(returnBytes))
	log.Debug(ctx, here, commandLog)
	return returnBytes, nil
}

//betterResponseName return response name but return ok when err=0
//
//	result := betterResponseName(errOK.XXX_MapID(), errOK)
func betterResponseName(id uint16, response interface{}) string {
	name := response.(Response).XXX_MapName()
	if id == 1 {
		err := response.(*shared.Err)
		if err.Code == 0 {
			return "OK"
		}
		return fmt.Sprintf("Err=%v", err.Code)
	}
	return name
}

// timeExecuteAction execute action and measure time, log warning if too slow
func (dp *Dispatch) timeExecuteAction(ctx context.Context, action interface{}) (uint16, interface{}, error) {
	timer := tools.NewTimer()
	timer.Start()
	responseID, response, err := dp.runAction(ctx, action)
	ms := int(timer.Stop())
	slow := app.IsSlow(ms)
	if slow > 0 {
		log.Warning(ctx, here, fmt.Sprintf("%v is slow, expected finish in %v ms but it took %v ms", action.(Action).XXX_MapName(), int(slow), int(ms)))
	}
	return responseID, response, err
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
//
//when id <= 100 use shared map, id > 100 use dispatch map
func (dp *Dispatch) protoFromBuffer(id uint16, bytes []byte) (interface{}, error) {
	var obj interface{}
	if id <= 1000 {
		shareMap := &shared.MapXXX{}
		obj = shareMap.NewObjectByID(id)
	} else {
		obj = dp.Map.NewObjectByID(id)
	}
	if obj == nil {
		return nil, errors.Errorf("failed to map id %v", id)
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

// runAction send action to handler and get response
//
//normally we don't return err here, because we can log err to database let programmer to fix it. just return error response and error id let user track the problem
//
//DeadlineExceeded is the only error return
func (dp *Dispatch) runAction(ctx context.Context, action interface{}) (uint16, interface{}, error) {
	responseInterface, err := action.(Action).Main(ctx)
	if err != nil {
		return 0, nil, err
	}
	if responseInterface == nil {
		return 0, nil, errors.New("failed to get response from action.Main()")
	}
	response := responseInterface.(Response)
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
	protoInterface, err := dp.protoFromBuffer(id, protoBytes)
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to convert buffer to proto")
	}
	return id, protoInterface, nil
}

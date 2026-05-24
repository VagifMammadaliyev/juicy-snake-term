package messages

import (
	"fmt"
	"io"
)

type MsgHandler interface {
	Handle(io.Writer, []byte) error
}

var handlersRegistry = make(map[MsgType]MsgHandler)

func RegisterHandler(msgType MsgType, handler MsgHandler) {
	handlersRegistry[msgType] = handler
}

func GetHandler(msgType MsgType) (MsgHandler, bool) {
	handler, ok := handlersRegistry[msgType]
	return handler, ok
}

func HandleMessage(w io.Writer, data []byte) error {
	if len(data) == 0 {
		return nil // no handling
	}

	msgType := MsgType(data[0])
	handler, ok := GetHandler(msgType)
	if !ok {
		return fmt.Errorf("handler not found for message type: %d", msgType)
	}

	return handler.Handle(w, data[1:])
}

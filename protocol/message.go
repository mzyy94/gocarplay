package protocol

import (
	"bytes"
	"errors"
	"io"
	"reflect"

	"github.com/lunixbochs/struc"
)

var messageTypes = map[reflect.Type]uint32{
	reflect.TypeOf(&SendFile{}):         0x99,
	reflect.TypeOf(&Open{}):             0x01,
	reflect.TypeOf(&Heartbeat{}):        0xaa,
	reflect.TypeOf(&ManufacturerInfo{}): 0x14,
	reflect.TypeOf(&CarPlay{}):          0x08,
	reflect.TypeOf(&SoftwareVersion{}):  0xcc,
	reflect.TypeOf(&BluetoothAddress{}): 0x0a,
	reflect.TypeOf(&BluetoothPIN{}):     0x0c,
	reflect.TypeOf(&Plugged{}):          0x02,
	reflect.TypeOf(&Unplugged{}):        0x04,
	reflect.TypeOf(&VideoData{}):        0x06,
	reflect.TypeOf(&AudioData{}):        0x07,
	reflect.TypeOf(&Touch{}):            0x05,
}

// Message is header structure of data protocol
type Message struct {
	Magic   uint32 `struc:"uint32,little"`
	Length  uint32 `struc:"uint32,little,sizeof=Payload"`
	Type    uint32 `struc:"uint32,little"`
	TypeN   uint32 `struc:"uint32,little"`
	Payload []byte
}

func PackMessage(buffer io.Writer, payload interface{}) error {
	var buf bytes.Buffer
	err := struc.Pack(&buf, payload)
	if err != nil {
		return err
	}
	msgType, found := messageTypes[reflect.TypeOf(payload)]
	if !found {
		return errors.New("No message found")
	}
	msgTypeN := (msgType ^ 0xffffffff) & 0xffffffff
	msg := &Message{Magic: 0x55aa55aa, Type: msgType, TypeN: msgTypeN, Payload: buf.Bytes()}
	return struc.Pack(buffer, msg)
}

func UnpackMessage(buffer io.Reader) (interface{}, error) {
	var msg Message
	err := struc.Unpack(buffer, &msg)
	if err != nil {
		return nil, err
	}

	for key, value := range messageTypes {
		if value == msg.Type {
			payload := reflect.New(key.Elem()).Interface()
			struc.Unpack(bytes.NewBuffer(msg.Payload), payload)
			return payload, nil
		}
	}

	return nil, errors.New("No message found")
}

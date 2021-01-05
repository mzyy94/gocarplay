package protocol

import (
	"bytes"
	"errors"
	"io"
	"reflect"

	"github.com/lunixbochs/struc"
)

const magicNumber uint32 = 0x55aa55aa

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
	if reflect.ValueOf(payload).Elem().NumField() > 0 {
		err := struc.Pack(&buf, payload)
		if err != nil {
			return err
		}
	}
	msgType, found := messageTypes[reflect.TypeOf(payload)]
	if !found {
		return errors.New("No message found")
	}
	msgTypeN := (msgType ^ 0xffffffff) & 0xffffffff
	msg := &Message{Magic: magicNumber, Type: msgType, TypeN: msgTypeN, Payload: buf.Bytes()}
	return struc.Pack(buffer, msg)
}

func UnpackHeader(buffer io.Reader, msg *Message) error {
	err := struc.Unpack(buffer, msg)
	if err != nil {
		return err
	}
	if msg.Magic != magicNumber {
		return errors.New("Invalid magic number")
	}
	if (msg.Type^0xffffffff)&0xffffffff != msg.TypeN {
		return errors.New("Invalid type")
	}
	return nil
}

func UnpackPayload(msgType uint32, buffer io.Reader) (interface{}, error) {
	for key, value := range messageTypes {
		if value == msgType {
			payload := reflect.New(key.Elem()).Interface()
			struc.Unpack(buffer, payload)
			return payload, nil
		}
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, buffer)
	return &Unknown{Data: buf.Bytes()}, nil
}

func UnpackMessage(buffer io.Reader) (interface{}, error) {
	var msg Message
	err := UnpackHeader(buffer, &msg)
	if err != nil {
		return nil, err
	}
	return UnpackPayload(msg.Type, bytes.NewBuffer(msg.Payload))
}

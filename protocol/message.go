package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"

	"github.com/lunixbochs/struc"
)

const magicNumber uint32 = 0x55aa55aa

var messageTypes = map[reflect.Type]uint32{
	reflect.TypeOf(&SendFile{}):            0x99,
	reflect.TypeOf(&Open{}):                0x01,
	reflect.TypeOf(&Heartbeat{}):           0xaa,
	reflect.TypeOf(&ManufacturerInfo{}):    0x14,
	reflect.TypeOf(&CarPlay{}):             0x08,
	reflect.TypeOf(&SoftwareVersion{}):     0xcc,
	reflect.TypeOf(&BluetoothAddress{}):    0x0a,
	reflect.TypeOf(&BluetoothPIN{}):        0x0c,
	reflect.TypeOf(&Plugged{}):             0x02,
	reflect.TypeOf(&Unplugged{}):           0x04,
	reflect.TypeOf(&VideoData{}):           0x06,
	reflect.TypeOf(&AudioData{}):           0x07,
	reflect.TypeOf(&Touch{}):               0x05,
	reflect.TypeOf(&BluetoothDeviceName{}): 0x0d,
	reflect.TypeOf(&WifiDeviceName{}):      0x0e,
	reflect.TypeOf(&BluetoothPairedList{}): 0x12,
}

// Header is header structure of data protocol
type Header struct {
	Magic  uint32 `struc:"uint32,little"`
	Length uint32 `struc:"uint32,little"`
	Type   uint32 `struc:"uint32,little"`
	TypeN  uint32 `struc:"uint32,little"`
}

func packPayload(buffer io.Writer, payload interface{}) error {
	if reflect.ValueOf(payload).Elem().NumField() > 0 {
		return struc.Pack(buffer, payload)
	}
	// Nothing to do
	return nil
}

func packHeader(payload interface{}, buffer io.Writer, data []byte) error {
	msgType, found := messageTypes[reflect.TypeOf(payload)]
	if !found {
		return errors.New("No message found")
	}
	msgTypeN := (msgType ^ 0xffffffff) & 0xffffffff
	msg := &Header{Magic: magicNumber, Length: uint32(len(data)), Type: msgType, TypeN: msgTypeN}
	err := struc.Pack(buffer, msg)
	if err != nil {
		return err
	}
	_, err = buffer.Write(data)
	return err
}

func Marshal(payload interface{}) ([]byte, error) {
	var buf, buffer bytes.Buffer
	err := packPayload(&buf, payload)
	if err != nil {
		return nil, err
	}
	err = packHeader(payload, &buffer, buf.Bytes())
	return buffer.Bytes(), err
}

func GetPayloadByHeader(hdr Header) interface{} {
	for key, value := range messageTypes {
		if value == hdr.Type {
			return reflect.New(key.Elem()).Interface()
		}
	}
	return &Unknown{Type: hdr.Type}
}

func Unmarshal(data []byte, payload interface{}) error {
	if len(data) > 0 {
		err := struc.Unpack(bytes.NewBuffer(data), payload)
		if err != nil {
			return err
		}
	}

	switch payload := payload.(type) {
	case *Header:
		if payload.Magic != magicNumber {
			return errors.New("Invalid magic number")
		}
		if (payload.Type^0xffffffff)&0xffffffff != payload.TypeN {
			return errors.New("Invalid type")
		}
	case *AudioData:
		switch len(data) - 12 {
		case 1:
			payload.Command = AudioCommand(data[12])
		case 4:
			binary.Read(bytes.NewBuffer(data[12:]), binary.LittleEndian, &payload.VolumeDuration)
		default:
			payload.Data = data[12:]
		}
	case *BluetoothDeviceName:
		payload.Data = string(data)
	case *WifiDeviceName:
		payload.Data = string(data)
	case *BluetoothPairedList:
		payload.Data = string(data)
	case *Unknown:
		payload.Data = data
	}

	return nil
}

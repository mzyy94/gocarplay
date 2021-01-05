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

// Header is header structure of data protocol
type Header struct {
	Magic  uint32 `struc:"uint32,little"`
	Length uint32 `struc:"uint32,little"`
	Type   uint32 `struc:"uint32,little"`
	TypeN  uint32 `struc:"uint32,little"`
}

func PackPayload(buffer io.Writer, payload interface{}) error {
	if reflect.ValueOf(payload).Elem().NumField() > 0 {
		return struc.Pack(buffer, payload)
	}
	// Nothing to do
	return nil
}

func PackHeader(payload interface{}, buffer io.Writer, data []byte) error {
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

func PackMessage(buffer io.Writer, payload interface{}) error {
	var buf bytes.Buffer
	err := PackPayload(&buf, payload)
	if err != nil {
		return err
	}
	return PackHeader(payload, buffer, buf.Bytes())
}

func UnpackHeader(buffer io.Reader, hdr *Header) error {
	err := struc.Unpack(buffer, hdr)
	if err != nil {
		return err
	}
	if hdr.Magic != magicNumber {
		return errors.New("Invalid magic number")
	}
	if (hdr.Type^0xffffffff)&0xffffffff != hdr.TypeN {
		return errors.New("Invalid type")
	}
	return nil
}

func UnpackPayload(hdr Header, buffer io.Reader) (interface{}, error) {
	for key, value := range messageTypes {
		if value == hdr.Type {
			payload := reflect.New(key.Elem()).Interface()
			struc.Unpack(buffer, payload)

			switch payload := payload.(type) {
			case *AudioData:
				buf := new(bytes.Buffer)
				io.Copy(buf, buffer)
				bin := buf.Bytes()

				switch len(bin) {
				case 1:
					payload.Command = AudioCommand(bin[0])
				case 4:
					binary.Read(bytes.NewBuffer(bin), binary.LittleEndian, &payload.VolumeDuration)
				default:
					payload.Data = bin
				}
			}
			reflect.ValueOf(payload).Elem().FieldByName("Header").Set(reflect.ValueOf(hdr))

			return payload, nil
		}
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, buffer)
	return &Unknown{Header: hdr, Data: buf.Bytes()}, nil
}

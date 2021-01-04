package protocol

import (
	"bytes"
	"errors"
	"reflect"

	"github.com/lunixbochs/struc"
)

// Message is header structure of data protocol
type Message struct {
	Magic   uint32 `struc:"uint32,little"`
	Length  uint32 `struc:"uint32,little,sizeof=Payload"`
	Type    uint32 `struc:"uint32,little"`
	TypeN   uint32 `struc:"uint32,little"`
	Payload []byte
}

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

type SendFile struct {
	FileNameSize int32 `struc:"int32,little,sizeof=FileName"`
	FileName     string
	ContentSize  int32 `struc:"int32,little,sizeof=Content"`
	Content      []byte
}

type Open struct {
	Width          int32 `struc:"int32,little"`
	Height         int32 `struc:"int32,little"`
	VideoFrameRate int32 `struc:"int32,little"`
	Format         int32 `struc:"int32,little"`
	PacketMax      int32 `struc:"int32,little"`
	IBoxVersion    int32 `struc:"int32,little"`
	PhoneWorkMode  int32 `struc:"int32,little"`
}

type Heartbeat struct{}

type ManufacturerInfo struct {
	A int32 `struc:"int32,little"`
	B int32 `struc:"int32,little"`
}

type CarPlay struct {
	Type CarPlayType `struc:"int32,little"`
}

type SoftwareVersion struct {
	Version string `struc:"[32]byte"`
}

type BluetoothAddress struct {
	Address string `struc:"[17]byte"`
}

type BluetoothPIN struct {
	Address string `struc:"[4]byte"`
}

type Plugged struct {
	PhoneType int  `struc:"int32,little"`
	Wifi      bool `struc:"int32,little"`
	// FIXME: Send WifiParam only when no wifi is ok
}

type Unplugged struct{}

type VideoData struct {
	Width    int32  `struc:"int32,little"`
	Height   int32  `struc:"int32,little"`
	Flags    int32  `struc:"int32,little"`
	Unknown1 int32  `struc:"int32,little"`
	Unknown2 int32  `struc:"int32,little"`
	Data     []byte `struc:"[]byte"`
}

type AudioData struct {
	DecodeType int32     `struc:"int32,little"`
	Volume     float32   `struc:"float32,little"`
	AudioType  AudioType `struc:"int32,little"`
	// Command AudioCommand  `struc:"int32,little"`
	// VolumeDuration int32  `struc:"int32,little"`
	Data []byte `struc:"[]byte"`
}

type Touch struct {
	Action TouchAction `struc:"int32,little"`
	X      int32       `struc:"int32,little"`
	Y      int32       `struc:"int32,little"`
	Flags  uint32      `struc:"int32,little"`
}

type MultiTouch struct {
	// TODO: Implement
}

///

func PackMessage(payload interface{}, buffer *bytes.Buffer) error {
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

func UnpackMessage(buffer *bytes.Buffer) (interface{}, error) {
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

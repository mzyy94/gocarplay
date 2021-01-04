package protocol

import (
	"bytes"
	"errors"
	"reflect"

	"github.com/lunixbochs/struc"
)

/// enum type def

type CarPlayType uint32

const (
	Invalid           = CarPlayType(0)
	BtnSiri           = CarPlayType(5)
	CarMicrophone     = CarPlayType(7)
	BtnLeft           = CarPlayType(100)
	BtnRight          = CarPlayType(101)
	BtnSelectDown     = CarPlayType(104)
	BtnSelectUp       = CarPlayType(105)
	BtnBack           = CarPlayType(106)
	BtnDown           = CarPlayType(114)
	BtnHome           = CarPlayType(200)
	BtnPlay           = CarPlayType(201)
	BtnPause          = CarPlayType(202)
	BtnNextTrack      = CarPlayType(204)
	BtnPrevTrack      = CarPlayType(205)
	SupportWifi       = CarPlayType(1000)
	SupportWifiNeedKo = CarPlayType(1012)
)

type AudioCommand uint32

const (
	AudioOutputStart    = AudioCommand(0x01)
	AudioOutputStop     = AudioCommand(0x02)
	AudioInputConfig    = AudioCommand(0x03)
	AudioPhonecallStart = AudioCommand(0x04)
	AudioPhonecallStop  = AudioCommand(0x05)
	AudioNaviStart      = AudioCommand(0x06)
	AudioNaviStop       = AudioCommand(0x07)
	AudioSiriStart      = AudioCommand(0x08)
	AudioSiriStop       = AudioCommand(0x09)
	AudioMediaStart     = AudioCommand(0x0a)
	AudioMediaStop      = AudioCommand(0x0b)
)

type AudioFormat struct {
	Frequency, Channel, Bitrate int
}

type AudioType uint32

var AudioDecodeTypes = map[AudioType]AudioFormat{
	0: {0, 0, 0},
	1: {44100, 2, 16},
	2: {44100, 2, 16},
	3: {8000, 1, 16},
	4: {48000, 2, 16},
	5: {16000, 1, 16},
	6: {24000, 1, 16},
	7: {16000, 2, 16},
}

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

type TouchAction uint32

const (
	TouchDown = TouchAction(14)
	TouchMove = TouchAction(15)
	TouchUp   = TouchAction(16)
)

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

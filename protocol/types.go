package protocol

import (
	"fmt"
	"strings"
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

func (c CarPlayType) GoString() string {
	switch c {
	case 0:
		return "Invalid"
	case 5:
		return "BtnSiri"
	case 7:
		return "CarMicrophone"
	case 100:
		return "BtnLeft"
	case 101:
		return "BtnRight"
	case 104:
		return "BtnSelectDown"
	case 105:
		return "BtnSelectUp"
	case 106:
		return "BtnBack"
	case 114:
		return "BtnDown"
	case 200:
		return "BtnHome"
	case 201:
		return "BtnPlay"
	case 202:
		return "BtnPause"
	case 204:
		return "BtnNextTrack"
	case 205:
		return "BtnPrevTrack"
	case 1000:
		return "SupportWifi"
	case 1012:
		return "SupportWifiNeedKo"
	}
	return fmt.Sprintf("Unknown(%d)", c)
}

type AudioCommand uint8

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

func (c AudioCommand) GoString() string {
	switch c {
	case 0x01:
		return "AudioOutputStart"
	case 0x02:
		return "AudioOutputStop"
	case 0x03:
		return "AudioInputConfig"
	case 0x04:
		return "AudioPhonecallStart"
	case 0x05:
		return "AudioPhonecallStop"
	case 0x06:
		return "AudioNaviStart"
	case 0x07:
		return "AudioNaviStop"
	case 0x08:
		return "AudioSiriStart"
	case 0x09:
		return "AudioSiriStop"
	case 0x0a:
		return "AudioMediaStart"
	case 0x0b:
		return "AudioMediaStop"
	}
	return fmt.Sprintf("Unknown(%d)", c)
}

type AudioFormat struct {
	Frequency, Channel, Bitrate uint16
}

type DecodeType uint32

var AudioDecodeTypes = map[DecodeType]AudioFormat{
	0: {0, 0, 0},
	1: {44100, 2, 16},
	2: {48000, 2, 16},
	3: {8000, 1, 16},
	4: {48000, 2, 16},
	5: {16000, 1, 16},
	6: {24000, 1, 16},
	7: {16000, 2, 16},
}

type TouchAction uint32

const (
	TouchDown = TouchAction(14)
	TouchMove = TouchAction(15)
	TouchUp   = TouchAction(16)
)

type NullTermString string

func (s NullTermString) GoString() string {
	return "'" + strings.TrimRight(string(s), "\x00") + "'"
}

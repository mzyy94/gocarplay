package protocol

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

type TouchAction uint32

const (
	TouchDown = TouchAction(14)
	TouchMove = TouchAction(15)
	TouchUp   = TouchAction(16)
)

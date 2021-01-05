package protocol

type SendFile struct {
	Header       `struc:"skip"`
	FileNameSize int32 `struc:"int32,little,sizeof=FileName"`
	FileName     string
	ContentSize  int32 `struc:"int32,little,sizeof=Content"`
	Content      []byte
}

type Open struct {
	Header         `struc:"skip"`
	Width          int32 `struc:"int32,little"`
	Height         int32 `struc:"int32,little"`
	VideoFrameRate int32 `struc:"int32,little"`
	Format         int32 `struc:"int32,little"`
	PacketMax      int32 `struc:"int32,little"`
	IBoxVersion    int32 `struc:"int32,little"`
	PhoneWorkMode  int32 `struc:"int32,little"`
}

type Heartbeat struct {
	Header `struc:"skip"`
}

type ManufacturerInfo struct {
	Header `struc:"skip"`
	A      int32 `struc:"int32,little"`
	B      int32 `struc:"int32,little"`
}

type CarPlay struct {
	Header `struc:"skip"`
	Type   CarPlayType `struc:"int32,little"`
}

type SoftwareVersion struct {
	Header  `struc:"skip"`
	Version string `struc:"[32]byte"`
}

type BluetoothAddress struct {
	Header  `struc:"skip"`
	Address string `struc:"[17]byte"`
}

type BluetoothPIN struct {
	Header  `struc:"skip"`
	Address string `struc:"[4]byte"`
}

type Plugged struct {
	Header    `struc:"skip"`
	PhoneType int  `struc:"int32,little"`
	Wifi      bool `struc:"int32,little"`
	// FIXME: Send WifiParam only when no wifi is ok
}

type Unplugged struct {
	Header `struc:"skip"`
}

type VideoData struct {
	Header   `struc:"skip"`
	Width    int32  `struc:"int32,little"`
	Height   int32  `struc:"int32,little"`
	Flags    int32  `struc:"int32,little"`
	Length   int32  `struc:"int32,little,sizeof=Data"`
	Unknown2 int32  `struc:"int32,little"`
	Data     []byte `struc:"[]byte"`
}

type AudioData struct {
	Header         `struc:"skip"`
	DecodeType     int32        `struc:"int32,little"`
	Volume         float32      `struc:"float32,little"`
	AudioType      AudioType    `struc:"int32,little"`
	Command        AudioCommand `struc:"skip"`
	VolumeDuration int32        `struc:"skip"`
	Data           []byte       `struc:"skip"`
}

type Touch struct {
	Header `struc:"skip"`
	Action TouchAction `struc:"int32,little"`
	X      uint32      `struc:"uint32,little"`
	Y      uint32      `struc:"uint32,little"`
	Flags  uint32      `struc:"int32,little"`
}

type MultiTouch struct {
	Header `struc:"skip"`
	// TODO: Implement
}

type Unknown struct {
	Header `struc:"skip"`
	Data   []byte
}

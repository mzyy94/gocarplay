package link

import (
	"bytes"
	"context"
	"encoding/binary"
	"time"

	"github.com/google/gousb"
	"github.com/mzyy94/gocarplay/protocol"
)

var epIn *gousb.InEndpoint
var epOut *gousb.OutEndpoint
var ctx context.Context
var Done func()

func Init() error {
	var err error
	epIn, epOut, Done, err = Connect()
	if err != nil {
		return err
	}
	ctx = context.Background()
	return nil
}

func intToByte(data int32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, data)
	return buf.Bytes()
}

func Start(width, height, fps, dpi int32) {
	SendData(&protocol.SendFile{FileName: "/tmp/screen_dpi\x00", Content: intToByte(dpi)})
	SendData(&protocol.Open{Width: width, Height: height, VideoFrameRate: fps, Format: 5, PacketMax: 4915200, IBoxVersion: 2, PhoneWorkMode: 2})

	SendData(&protocol.ManufacturerInfo{A: 0, B: 0})
	SendData(&protocol.SendFile{FileName: "/tmp/night_mode\x00", Content: intToByte(1)})
	SendData(&protocol.SendFile{FileName: "/tmp/hand_drive_mode\x00", Content: intToByte(1)})
	SendData(&protocol.SendFile{FileName: "/tmp/charge_mode\x00", Content: intToByte(0)})
	SendData(&protocol.SendFile{FileName: "/tmp/box_name\x00", Content: bytes.NewBufferString("BoxName").Bytes()})

	for {
		SendData(&protocol.Heartbeat{})
		time.Sleep(2 * time.Second)
	}
}

func Communicate(onData func(interface{}), onError func(error)) {
	for {
		received, err := ReceiveMessage(epIn, ctx)
		if err != nil {
			onError(err)
		} else {
			onData(received)
		}
	}
}

func SendData(data interface{}) error {
	return SendMessage(epOut, data)
}

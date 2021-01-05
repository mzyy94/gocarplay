package link

import (
	"bytes"

	"github.com/google/gousb"
	"github.com/mzyy94/gocarplay/protocol"
)

func SendMessage(epOut *gousb.OutEndpoint, msg interface{}) error {
	var buf bytes.Buffer
	err := protocol.PackMessage(&buf, msg)
	if err != nil {
		return err
	}
	_, err = epOut.Write(buf.Bytes()[:16])
	if err != nil {
		return err
	}
	if len(buf.Bytes()) > 16 {
		_, err = epOut.Write(buf.Bytes()[16:])
	}
	return err
}

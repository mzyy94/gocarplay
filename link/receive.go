package link

import (
	"bytes"
	"context"
	"io"

	"github.com/google/gousb"
	"github.com/mzyy94/gocarplay/protocol"
)

func ReceiveMessage(epIn *gousb.InEndpoint, ctx context.Context) (interface{}, error) {
	buf := make([]byte, 16)
	var msg protocol.Message
	num, err := epIn.ReadContext(ctx, buf)
	if err != nil {
		return nil, err
	}
	err = protocol.UnpackHeader(bytes.NewBuffer(buf[:num]), &msg)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if msg.Length > 0 {
		buf := make([]byte, msg.Length)
		num, err = epIn.ReadContext(ctx, buf)
		if err != nil {
			return nil, err
		}
		return protocol.UnpackPayload(msg.Type, bytes.NewBuffer(buf[:num]))
	}
	return protocol.UnpackPayload(msg.Type, &bytes.Buffer{})
}

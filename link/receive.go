package link

import (
	"context"

	"github.com/google/gousb"
	"github.com/mzyy94/gocarplay/protocol"
)

func ReceiveMessage(epIn *gousb.InEndpoint, ctx context.Context) (interface{}, error) {
	buf := make([]byte, 16)
	var hdr protocol.Header
	num, err := epIn.ReadContext(ctx, buf)
	if err != nil {
		return nil, err
	}
	err = protocol.Unmarshal(buf[:num], &hdr)
	if err != nil {
		return nil, err
	}

	payload := protocol.GetPayloadByHeader(hdr)
	buf = make([]byte, hdr.Length)

	if hdr.Length > 0 {
		num, err = epIn.ReadContext(ctx, buf)
		if err != nil {
			return nil, err
		}
	} else {
		num = 0
	}
	err = protocol.Unmarshal(buf[:num], payload)
	return payload, err
}

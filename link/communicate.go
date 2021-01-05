package link

import (
	"context"

	"github.com/google/gousb"
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

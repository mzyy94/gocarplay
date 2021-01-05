package link

import (
	"errors"
	"time"

	"github.com/google/gousb"
)

func Connect() (*gousb.InEndpoint, *gousb.OutEndpoint, func(), error) {
	cleanTask := make([]func(), 0)
	defer func() {
		for _, task := range cleanTask {
			task()
		}
	}()

	ctx := gousb.NewContext()
	cleanTask = append(cleanTask, func() { ctx.Close() })

	var (
		dev       *gousb.Device
		err       error
		waitCount = 5
	)

	for {
		dev, err = ctx.OpenDeviceWithVIDPID(0x1314, 0x1520)
		if err != nil {
			return nil, nil, nil, err
		}
		if dev == nil {
			waitCount--
			if waitCount < 0 {
				return nil, nil, nil, errors.New("Could not find a device")
			}
			time.Sleep(3 * time.Second)
			continue
		}
		cleanTask = append(cleanTask, func() { dev.Close() })
		break
	}

	intf, done, err := dev.DefaultInterface()
	if err != nil {
		return nil, nil, nil, err
	}
	cleanTask = append(cleanTask, done)

	epOut, err := intf.OutEndpoint(1)
	if err != nil {
		return nil, nil, nil, err
	}
	epIn, err := intf.InEndpoint(1)
	if err != nil {
		return nil, nil, nil, err
	}

	closeTask := make([]func(), len(cleanTask))
	copy(closeTask, cleanTask)
	cleanTask = nil

	return epIn, epOut, func() {
		for _, task := range closeTask {
			task()
		}
	}, nil
}

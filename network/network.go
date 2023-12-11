package network

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
	"github.com/running910/gokit/logger"
)

func Ping(dst string, src string) error {
	pinger, err := ping.NewPinger(dst)
	if err != nil {
		logger.Errorf("ping.NewPinger() failed! reason:%s", err)
		return err
	}

	pinger.Interval = time.Millisecond * 50
	pinger.Count = 2
	pinger.Timeout = time.Second * 2
	pinger.Source = src
	pinger.SetPrivileged(true)

	err = pinger.Run()
	if err != nil {
		logger.Debugf("ping.Run() failed! reason:%s", err)
		return err
	}

	stats := pinger.Statistics()

	if stats.PacketsRecv > 0 {
		return nil
	} else {
		logger.Debug("ping", dst, "failed with src", src)
		return fmt.Errorf("none packets recieved!")
	}
}

package domain

import (
	"net"
	"time"
)

type Pinger interface {
	Ping() (*PingResult, error)
	PingContinuous() error
	Traceroute() error
}

type PingResult struct {
	Peer      net.Addr
	TimeTaken time.Duration
}

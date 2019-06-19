package domain

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func DefaultSettings() []PingSettings {
	return []PingSettings{
		PingSettings{
			IPType:        "",
			SourceAddress: "::",
			TargetAddress: "2001:4860:4860::8888",
			Count:         10,
		},
		PingSettings{
			IPType:        "",
			SourceAddress: "192.168.0.108",
			TargetAddress: "8.8.8.8",
			Count:         10,
		},
	}
}

type PingSettings struct {
	IPType        string
	SourceAddress string
	TargetAddress string
	Count         int
}

func IPProtocol(address string) int {
	for i := 0; i < len(address); i++ {
		switch address[i] {
		case '.':
			return 4
		case ':':
			return 6
		}
	}
	panic(fmt.Errorf("invalid address '%v'. Neither dot nor : included", address))
}

func (settings *PingSettings) ProtocolFromSourceAddress() string {
	address := settings.SourceAddress

	for i := 0; i < len(address); i++ {
		switch address[i] {
		case '.':
			//return "ip4:imp"
			return "udp4"
		case ':':
			return "udp6"
		}
	}
	log.Fatalf("invalid address '%v'. Neither dot nor : included", address)
	return "udp6"
}

func (settings *PingSettings) MessageToTargetAddress() icmp.Message {
	proto := IPProtocol(settings.TargetAddress)
	switch proto {
	case 4:
		return icmp.Message{
			Type: ipv4.ICMPTypeEcho, Code: 0,
			Body: &icmp.Echo{
				ID: os.Getpid() & 0xffff, Seq: 1,
				Data: []byte("HELLO-R-U-THERE"),
			},
		}
	case 6:
		return icmp.Message{
			Type: ipv6.ICMPTypeEchoRequest, Code: 0,
			Body: &icmp.Echo{
				ID: os.Getpid() & 0xffff, Seq: 1,
				Data: []byte("HELLO-R-U-THERE"),
			},
		}
	}
	panic(fmt.Sprintf("invalid protocol %d", proto))
}

type Message struct {
	Value string
}

func (m *Message) Len(proto int) int {
	return len([]byte(m.Value))
}

func (m *Message) Marshal(proto int) ([]byte, error) {
	return []byte(m.Value), nil
}

func (settings *PingSettings) MessageToTargetAddressWithTTL(message Message) icmp.Message {
	proto := IPProtocol(settings.TargetAddress)
	switch proto {
	case 4:
		return icmp.Message{
			Type: ipv4.ICMPTypeEcho, Code: 0,
			Body: &message,
		}
	case 6:
		return icmp.Message{
			Type: ipv6.ICMPTypeEchoRequest, Code: 0,
			Body: &message,
		}
	}
	panic(fmt.Sprintf("invalid protocol %d", proto))
}
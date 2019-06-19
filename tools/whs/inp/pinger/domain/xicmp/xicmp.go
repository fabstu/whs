package xicmp

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/bradfitz/iter"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"

	"aduu.dev/tools/whs/inp/pinger/domain"
)

const (
	maxTTL = 128
)

func Make(settings domain.PingSettings) domain.Pinger {
	return &XICMPPing{
		PingSettings: settings,
	}
}

type XICMPPing struct {
	domain.PingSettings
}

// Get preferred outbound ip of this machine
func GetOutboundIP(target string) net.IP {
	proto := domain.IPProtocol(target)

	var testTarget string
	var targetProto string

	switch proto {
	case 4:
		testTarget = "8.8.8.8"
		targetProto = "ip4:icmp"
	case 6:
		testTarget = "[2001:4860:4860:0:0:0:0:8888]:80"
		targetProto = "udp"
	default:
		panic(fmt.Errorf("unknown ip protocol %d", proto))
	}
	fmt.Println("Test target:", testTarget)

	conn, err := net.Dial(targetProto, testTarget)

	if err != nil {
		panic(err)
	}
	defer conn.Close()


	var localAddr net.IP
	switch proto {
	case 4:
		localAddr = conn.LocalAddr().(*net.IPAddr).IP
	case 6:
		localAddr = conn.LocalAddr().(*net.UDPAddr).IP
	}

	return localAddr
}

func localAddresses() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
			continue
		}

		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPNet:
				fmt.Printf("%v : %s (%s)\n", i.Name, v, v.IP.DefaultMask())
			}

		}
	}
}

func protocolFromAddress(address string) string {
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

func (settings *XICMPPing) sendMessage(msg icmp.Message, conn *icmp.PacketConn) error {
	wb, err := msg.Marshal(nil)
	if err != nil {
		return err
	}

	target := &net.UDPAddr{
		IP:   net.ParseIP(settings.TargetAddress),
		Port: 0,
	}

	if _, err := conn.WriteTo(wb, target); err != nil {
		return err
	}
	return nil
}

func (settings *XICMPPing) sendEcho(conn *icmp.PacketConn) error {
	wm := settings.MessageToTargetAddress()
	return settings.sendMessage(wm, conn)
}

func (settings *XICMPPing) parsePingMessage(rb []byte, n int) (err error) {
	proto := domain.IPProtocol(settings.TargetAddress)
	var rm *icmp.Message
	fmt.Println("Proto:", proto)
	switch proto {
	case 4:
		rm, err = icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), rb[:n])
		if err != nil {
			return err
		}
	case 6:
		rm, err = icmp.ParseMessage(58, rb[:n])
		if err != nil {
			return err
		}


	default:
		panic(fmt.Sprintf("unknown proto %d", proto))
	}


	switch rm.Type {
	case ipv6.ICMPTypeEchoReply, ipv4.ICMPTypeEchoReply, ipv4.ICMPTypeExtendedEchoRequest:
		return nil
	default:
		return fmt.Errorf("not an echo response: %#v", rm)
	}
}

func setIPv6HopLimit(fd int, v int) error {
	err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IPV6, syscall.IPV6_UNICAST_HOPS, v)
	if err != nil {
		return os.NewSyscallError("setsockopt", err)
	}
	return nil
}

func (settings *XICMPPing) setTTL(conn *icmp.PacketConn, ttl int) error {
	// Set Time To Live.
	proto := domain.IPProtocol(settings.TargetAddress)
	switch proto {
	case 4:
		if err := conn.IPv4PacketConn().SetTTL(ttl); err != nil {
			return err
		}
	case 6:
		if err := conn.IPv6PacketConn().SetHopLimit(ttl); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown protocol %v", proto)
	}
	return nil
}

func listenToPackets(settings *XICMPPing, chosenSourceAddress string) (*icmp.PacketConn, error) {
	c, err := icmp.ListenPacket(
		protocolFromAddress(settings.SourceAddress), chosenSourceAddress)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (settings *XICMPPing) Ping() (*domain.PingResult, error) {
	chosenSourceAddress := GetOutboundIP(settings.TargetAddress).String()

	fmt.Println("chosen source:", chosenSourceAddress)

	c, err := listenToPackets(settings, chosenSourceAddress)
	if err != nil {
		return nil, err
	}

	if err := settings.setTTL(c, 128); err != nil {
		return nil, err
	}


	defer c.Close()

	rb := make([]byte, 500)
	start := time.Now()

	if err := settings.sendEcho(c); err != nil {
		return nil, err
	}

	n, peer, err := c.ReadFrom(rb)
	t := time.Now()
	if err != nil {
		return nil, err
	}

	if err := settings.parsePingMessage(rb, n); err != nil {
		return nil, err
	}

	return &domain.PingResult{
		Peer:      peer,
		TimeTaken: t.Sub(start),
	}, nil
}

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func (settings *XICMPPing) PingContinuous() error {
	for i := range iter.N(10) {
		fmt.Printf("i: %d", i)
		res, err := settings.Ping()
		if err != nil {
			fmt.Printf("ping failed: %v\n", err)
		} else {
			log.Printf("got reflection from %v in %v", res.Peer, res.TimeTaken)
		}
		time.Sleep(time.Second)
		fmt.Println()
	}
	return nil
}

/*
func (settings *XICMPPing) Traceroute() error {
	chosenSourceAddress := GetOutboundIP(settings.TargetAddress).String()

	fmt.Println("chosen source:", chosenSourceAddress)

	c, err := icmp.ListenPacket(
		protocolFromAddress(settings.SourceAddress), chosenSourceAddress)
	if err != nil {
		return err
	}
	defer c.Close()

	rb := make([]byte, 500)

	for ttl:= 0; ttl < maxTTL; ttl++ {
		start := time.Now()

		if err := settings.setTTL(c, ttl); err != nil {
			return fmt.Errorf("failed to set ttl: %v", err)
		}

		if err := settings.sendEcho(c, ttl); err != nil {
			return fmt.Errorf("failed to send message: %v", err)
		}

		n, peer, err := c.ReadFrom(rb)
		t := time.Now()
		if err != nil {
			return fmt.Errorf("failed to read: %v", err)
		}

		// TODO: Add more message types. Like TTL exceeded.
		//  And match senders.
		//  	Do I get a TTL back?
		//  	Do I use different source ports?
		if err := settings.parsePingMessage(rb, n); err != nil {
			return nil, err
		}
	}

}
*/
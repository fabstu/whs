package xicmp

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/icmp"

	"aduu.dev/tools/whs/inp/pinger/domain"
)

const (
	ProtocolICMP = 1
	ProtocolIPv6ICMP = 58
)

func (settings *XICMPPing) Traceroute() error {
	chosenSourceAddress := GetOutboundIP(settings.TargetAddress).String()

	fmt.Println("chosen source:", chosenSourceAddress)

	c, err := listenToPackets(settings, chosenSourceAddress)
	if err != nil {
		return  err
	}
	defer c.Close()

	m := make(map[int]string)

	go settings.readTraceMessages(m, c)

	for i := 0; i < 16; i++ {
		//fmt.Println("ttl:", i)

		if err := settings.setTTL(c, i); err != nil {
			return err
		}

		// start := time.Now()

		mymsg := domain.Message{
			Value: fmt.Sprintf("---%d", i),
		}

		fmt.Println("value:", mymsg.Value)

			msg := settings.MessageToTargetAddressWithTTL(mymsg)

		//fmt.Println("sending")
		if err := settings.sendMessage(msg, c); err != nil {
			return err
		}

		/*

		return &domain.PingResult{
			Peer:      peer,
			TimeTaken: t.Sub(start),
		}, nil
		*/
	}

	fmt.Println("finished sending. receiving.")
	time.Sleep(time.Second * 5)

	type pair struct {
		ttl int
		peer string
	}

	fmt.Println("result:")

	var pairs []pair
	for ttl, peer := range m {
		pairs = append(pairs, pair{
			ttl:  ttl,
			peer: peer,
		})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].ttl <= pairs[j].ttl
	})

	for _, pair := range pairs {
		fmt.Println(pair)
	}

	fmt.Println("exiting.")

	return nil
}

func (settings *XICMPPing) readTraceMessages(m map[int]string, c *icmp.PacketConn) {
	fmt.Println("reading")

	rb := make([]byte, 2000)
	for true {
		n, peer, err := c.ReadFrom(rb)
		//t := time.Now()
		if err != nil {
			return
		}

		fmt.Println("--------- n, peer:", n, peer)

		msg, err := settings.parseTraceMessage(rb, n)
		if err != nil {
			fmt.Println("error parsing:", err)
			continue
		}

		if msg == "" {
			fmt.Println("echo received")
			continue
		}

		split := strings.Split(msg, "-")
		if len(split) <= 1 {
			fmt.Println("got message which has no ttl.")
			continue
		}

		ttlString := split[len(split) - 1]

		fmt.Println("ttl-string:", ttlString)

		ttl, err := strconv.Atoi(strings.TrimSpace(ttlString))
		if err != nil {
			fmt.Println("failed to parse ttlString", ttlString)
			continue
		}

		fmt.Printf("ttl = %s\n", peer.String())
		m[ttl] = peer.String()
	}
}


func (settings *XICMPPing) parseTraceMessage(rb []byte, n int) (string, error) {
	proto := domain.IPProtocol(settings.TargetAddress)
	fmt.Println("Proto:", proto)

	msg, err := icmp.ParseMessage(ProtocolICMP, rb[0:n])
	if err != nil {
		fmt.Println("failed to parse message.")
		return "", err
	}
	fmt.Println("type:", msg.Type)

	switch v := msg.Body.(type) {
	case *icmp.TimeExceeded:
		fmt.Println("time exceeded.")
		fmt.Println("data:", string(v.Data))
		return string(v.Data), nil
	case *icmp.Echo:
		fmt.Println("is echo.")
		return "", nil
	default:
		return "", fmt.Errorf("unknown type: %v", msg.Type)
	}

	//fmt.Println("type:", msg.Type, "body:", msg.Body.(*domain.Message))
}

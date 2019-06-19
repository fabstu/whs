package dns

import (
	"fmt"

	"github.com/spf13/cobra"

	"aduu.dev/tools/aduu/helper"

	"github.com/miekg/dns"
)

const rootServer = "198.41.0.4" // A.ROOT-SERVER.NET

/*
Reverse DNS resolution. From NameserverIP to Nameserver-DNS-Name.
 */
func DomainForIP(m map[string]string, ip string) string {
	for nsName, nsIP := range m {
		if nsIP == ip {
			return nsName
		}
	}
	panic(fmt.Sprintf("failed to reverse-dns %s", ip))
}

func printMessage(msg *dns.Msg) {
	helper.PrintAny(msg)
}

func printARecord(a *dns.A) {
	helper.PrintAny(a)
}

/*
Resolves the given domain with the given nameserver-ip.
m is a nameserver -> ip-mapping.
*/
func resolve(domain string, nameserver string, m map[string]string) ([]string, error) {
	fmt.Printf("Choosing %s[%s] to resolve %s\n", DomainForIP(m, nameserver), nameserver, domain)

	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = false
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{Name: dns.Fqdn(domain), Qtype: dns.TypeA, Qclass: dns.ClassINET}

	c := new(dns.Client)
	in, rtt, err := c.Exchange(m1, nameserver+ ":53")
	if err != nil {
		return nil, fmt.Errorf("failed to ask root server: %v", err)
	}

	//fmt.Println("rtt:", rtt.String())
	_ = rtt
	//printMessage(m1)
	//printMessage(in)

	if in.Answer != nil {
		var out []string
		for _, answer := range in.Answer {
			switch a := answer.(type) {
			case *dns.A:
				//fmt.Println("A-Record:")
				out = append(out, a.A.String())
				//PrintAny(a)
			default:
				fmt.Println("Non-A Record as answer:")
				helper.PrintAny(a)
			}
		}
		fmt.Printf("Answers for %s: %v\n", domain, out)
		return out, nil
	}
	nextNameserver := ""
	out:
	for _, rr := range in.Ns {
		switch ns := rr.(type) {
		case *dns.NS:
			nextNameserver = ns.Ns
			break out
		case *dns.SOA:
			helper.PrintAny(rr)
		default:
			fmt.Printf("unknown .Ns rr")
			helper.PrintAny(rr)
		}
	}
	if nextNameserver == "" {
		return nil, fmt.Errorf("no Answer ans no Nameserver")
	}

	// Add to my map.
	for _, e := range in.Extra {
		switch a := e.(type) {
		case *dns.A:
			m[a.Hdr.Name] = a.A.String()
		default:
			continue
		}
	}

	var nameserverIP string

	// Resolving nameserver with extras.
	if nsIP, ok := m[nextNameserver]; ok {
		nameserverIP = nsIP
	} else {
		// Resolve nameserver itself.
		ips, err := resolveWithMap(nextNameserver, m)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve nameserver: %v", err)
		}
		if len(ips) == 0 {
			return nil, fmt.Errorf("no ip found for nameserver %s", nextNameserver)
		}

		nameserverIP = ips[0]
		m[nextNameserver] = nameserverIP
	}

	//fmt.Println("m:", m)

	return resolve(domain, nameserverIP, m)
}

/*
Resolve resolves the given domain-name recursively for you and returns either the resulting ip-addresses or a non-nil error.
 */
func Resolve(domain string) ([]string, error) {
	m := make(map[string]string)
	m["root"] = rootServer
	return resolveWithMap(domain, m)
}

func resolveWithMap(domain string, m map[string]string) ([]string, error) {
	fmt.Println("Resolving", domain)
	return resolve(domain, rootServer, m)
}

var resolveCMD = &cobra.Command{
	Use:   "resolve <domain-name>",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		//domain := "smtp.kolabnow.com"
		domain := args[0]

		res, err := Resolve(domain)
		if err != nil {
			return err
		}

		fmt.Println("Answers:")
		for _, a := range res {
			fmt.Println(a)
		}

		return nil
	},
	Args: cobra.ExactArgs(1),
}

func init() {

}
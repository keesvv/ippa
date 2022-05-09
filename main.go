package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type AddrFormatter []net.Addr

func (a AddrFormatter) String() string {
	var rawAddr []string

	for _, addr := range a {
		var flags []string

		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			log.Fatalln(err)
		}

		if ip.To4() != nil {
			flags = append(flags, "inet")
		} else {
			flags = append(flags, "inet6")
		}

		if ip.IsPrivate() {
			flags = append(flags, "private")
		}

		if ip.IsLoopback() {
			flags = append(flags, "lo")
		}

		rawAddr = append(rawAddr, fmt.Sprintf(
			"  %s %s", ip.String(),
			strings.Join(flags, " "),
		))
	}

	return strings.Join(rawAddr, "\n")
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalln(err)
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%s\n%s\n", iface.Name, AddrFormatter(addrs))
	}
}

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type AddrFormatter struct{ ip net.IP }

func (addr AddrFormatter) String() string {
	flags := []string{addr.ip.String()}

	if addr.ip.To4() != nil {
		flags = append(flags, "inet")
	} else {
		flags = append(flags, "inet6")
	}

	if addr.ip.IsPrivate() {
		flags = append(flags, "private")
	}

	if addr.ip.IsLoopback() {
		flags = append(flags, "lo")
	}

	return strings.Join(flags, " ")
}

func filterIfaces(filter []string, ifaces []net.Interface) []net.Interface {
	res := make([]net.Interface, 0)

	for _, f := range filter {
		for _, iface := range ifaces {
			if iface.Name == f {
				res = append(res, iface)
			}
		}
	}

	return res
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	filter := os.Args[1:]
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalln(err)
	}

	if len(filter) > 0 {
		ifaces = filterIfaces(filter, ifaces)
	}

	for _, iface := range ifaces {
		fmt.Println(iface.Name)

		addrs, err := iface.Addrs()
		if err != nil {
			log.Println(err)
			continue
		}

		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				log.Println(err)
				continue
			}

			fmt.Printf("  %s\n", AddrFormatter{ip})
		}
	}
}

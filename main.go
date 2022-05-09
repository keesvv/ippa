package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type AddrFormatter struct{ ip net.IP }

type AddrFlag struct {
	name    string
	include bool
}

func (addr AddrFormatter) String() string {
	var flags []string

	// not using map[string]bool here, maps
	// are not the right data structure for
	// this as maps do not indicate order
	flagMap := []AddrFlag{
		{"inet", addr.ip.To4() != nil},
		{"inet6", addr.ip.To4() == nil},
		{"private", addr.ip.IsPrivate()},
		{"lo", addr.ip.IsLoopback()},
	}

	for _, flag := range flagMap {
		if flag.include {
			flags = append(flags, flag.name)
		}
	}

	return strings.Join(append(
		[]string{addr.ip.String()}, flags...,
	), " ")
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

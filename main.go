package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error, please supply domain name as first argument")
		os.Exit(1)
	}
	domainName := os.Args[1]

	nameservers := map[string]string{
		"Google": "8.8.8.8:53",
		//"Google IPv6":     "[2001:4860:4860::8888]:53",
		"Cloudflare": "1.1.1.1:53",
		//"Cloudflare IPv6": "[2606:4700:4700::1111]:53",
		"OpenDNS": "208.67.222.222:53",
		//"OpenDNS IPv6":    "[2620:119:35::35]:53",
		"Freifunk MUC":         "5.1.66.255:53",   //IPv6: 2001:678:e68:f000::"
		"Censurfridns Denmark": "89.233.43.71:53", //IPv6: 2001:67c:28a4::
	}

	for nsName, nsIP := range nameservers {
		resolveAndPrint(nsName, nsIP, domainName)
	}

	fmt.Printf("Resolving %s with %20s: ", domainName, "native server")
	ipsNative, err := net.LookupHost(domainName)
	if err != nil {
		fmt.Println(strings.Split(err.Error(), ": ")[1])
		return
	}
	for _, ip := range ipsNative {
		fmt.Print(ip + " ")
	}
	fmt.Println()
}

func resolveAndPrint(nsName, nsIP, domainName string) {
	fmt.Printf("Resolving %s with %20s: ", domainName, nsName)
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, nsIP)
		},
	}
	resolvedIPs, err := r.LookupHost(context.Background(), domainName)
	if err != nil {
		fmt.Println(strings.Split(err.Error(), ": ")[1])
		return
	}
	for _, ip := range resolvedIPs {
		fmt.Print(ip + " ")
	}
	fmt.Println()
}

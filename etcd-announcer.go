/* A general purpose etcd service-announcer.

This application announces a service in etcd.

Examples:

1.
etcd-announce -service riak -msg foo

Would write a key "/service/riak" with content "foo".

2.
etcd-announce -path /service/riak -service riak-1 -type net -if eth0

Would announce the IP address of interface "eth0" in the key
"/service/riak/riak-1"

By default all announcements have a TTL of 30 seconds. The announcer
runs every TTL/2 seconds.

*/
package main

import (
	"flag"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"net"
	"os"
	"strings"
	"time"
)

func printUsageAndExit() {
	fmt.Println(`
etcd-announcer -service SERVICE

Other arguments:
-type <net|msg> (default: msg)
-if INTERFACE_NAME (must be set when -type net)
-msg MESSAGE (must be set when -type msg)
-ttl <TTL in seconds> (default: 30s)
-path ETCD_PATH (default: /service)
`)
	os.Exit(1)
}

func announce(key string, value string, etcdAddr string, ttl uint64) {
	etcdClient := etcd.NewClient([]string{etcdAddr})
	resp, err := etcdClient.Set(key, value, ttl)
	if err != nil {
		fmt.Println("An error occurred: ", err)
		os.Exit(1)
	}

	fmt.Printf("Announced %v: %v\n", resp.Node.Key, resp.Node.Value)
}

func startAnnouncer(ttl uint64, announcer func()) {
	var tickLength int64 = int64(ttl) / 2
	for _ = range time.Tick(time.Duration(tickLength) * time.Second) {
		announcer()
	}
}

// Returns the first IP address of the supplied network
// interface as a string. Fails if the interface does
// not exist or has no attached IP.
func getLocalIp(interfaceName string) string {
	localIf, err := net.InterfaceByName(interfaceName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	localIps, err := localIf.Addrs()
	if (err != nil) || (len(localIps) == 0) {
		fmt.Println("No IP for interface", interfaceName)
		os.Exit(1)
	}

	firstIp := localIps[0].String()
	return strings.Split(firstIp, "/")[0]
}

// Parse arguments and call main loop
func main() {
	service := flag.String("service", "", "Service name to announce")
	etcdAddr := flag.String("etcd", "http://127.0.0.1:4001",
		"etcd peer address. Defaults to http://127.0.0.1:4001")
	etcdPath := flag.String("path", "/service",
		"Path in etcd to place announcement in.")
	etcdTtl := flag.Uint64("ttl", 30, "Service announce TTL (default 30s)")
	announceType := flag.String("type", "msg",
		"Announce type. One of \"net\" or \"msg\"")
	announceMsg := flag.String("msg", "", "Message to announce")
	announceInterface := flag.String("if", "", "Interface to announce")

	flag.Parse()

	// Do some simple validation
	if *service == "" {
		fmt.Println("Must supply service name with -service")
		printUsageAndExit()
	}

	if (*announceType == "msg") && (*announceMsg == "") {
		fmt.Println("Must set message with -msg")
		printUsageAndExit()
	}

	if (*announceType == "net") && (*announceInterface == "") {
		fmt.Println("Must set interface with -if")
		printUsageAndExit()
	}

	// Set up announcer
	var announcer func()

	if *announceType == "msg" {
		announcer = func() {
			announce((*etcdPath + "/" + *service),
				*announceMsg, *etcdAddr, *etcdTtl)
		}
	} else if *announceType == "net" {
		announcer = func() {
			announce((*etcdPath + "/" + *service),
				getLocalIp(*announceInterface),
				*etcdAddr, *etcdTtl)
		}
	}

	// Announce once and start main loop
	announcer()

	go startAnnouncer(*etcdTtl, announcer)
	select {}
}

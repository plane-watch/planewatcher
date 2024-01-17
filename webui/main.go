package main

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

// returns a list of network devices (excludes loopback and bridge etc)
func getNetworkDevices() ([]netlink.Link, error) {

	// get list if network interfaces
	ll, err := netlink.LinkList()
	if err != nil {
		return []netlink.Link{}, err
	}

	// remove loopback
	llDevicesOnly := []netlink.Link{}
	for _, l := range ll {
		if l.Type() == "device" {
			if l.Attrs().Flags&net.FlagLoopback != net.FlagLoopback { // bitwise AND to check for loopback
				llDevicesOnly = append(llDevicesOnly, l)
			}
		}
	}

	// return sanitised list
	return llDevicesOnly, nil
}

func main() {

	ll, err := getNetworkDevices()
	if err != nil {
		panic(err)
	}

	for _, l := range ll {
		fmt.Println("attrs", *l.Attrs())
		fmt.Println(fmt.Sprintf("type %s", l.Type()))
	}
}

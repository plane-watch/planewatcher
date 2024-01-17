package main

import (
	"fmt"

	"github.com/plane-watch/planewatcher/webui/lib/network"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func main() {

	ll, err := network.GetNetworkDevices()
	if err != nil {
		panic(err)
	}

	for _, l := range ll {
		ips, err := netlink.AddrList(l, unix.AF_INET)
		if err != nil {
			panic(err)
		}
		fmt.Println(
			l.Attrs().Index,
			l.Attrs().Name,
			l.Attrs().OperState.String(),
			ips,
		)
		// fmt.Println("attrs", *l.Attrs())
		// fmt.Println(fmt.Sprintf("type %s", l.Type()))
	}
}

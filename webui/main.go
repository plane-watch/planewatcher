package main

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

func main() {

	ll, err := netlink.LinkList()
	if err != nil {
		panic(err)
	}

	fmt.Println(ll)

}

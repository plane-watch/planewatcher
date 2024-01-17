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

	for _, l := range ll {
		fmt.Println("attrs", *l.Attrs())
		fmt.Println(fmt.Sprintf("type %s", l.Type()))
	}
}

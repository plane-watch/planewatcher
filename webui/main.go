package main

import (
	"fmt"

	"github.com/plane-watch/planewatcher/webui/lib/network"
)

func main() {

	ll, err := network.GetNetworkDevices()
	if err != nil {
		panic(err)
	}

	for _, l := range ll {
		fmt.Println(
			l.Attrs().Index,
			l.Attrs().Name,
			l.Attrs().OperState,
		)
		// fmt.Println("attrs", *l.Attrs())
		// fmt.Println(fmt.Sprintf("type %s", l.Type()))
	}
}

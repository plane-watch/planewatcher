package main

import (
	"firstboot/lib/netplan"
	"fmt"
)

func main() {

	np, err := netplan.Load("/etc/netplan/planewatcher.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(np)

}

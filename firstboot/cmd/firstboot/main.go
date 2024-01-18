package main

import (
	"firstboot/lib/netplan"
	"fmt"
	"os"
)

const (
	netplanFile = "/etc/netplan/planewatcher.yaml"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {

	// check if netplan file exists
	// if ! fileExists(netplanFile) {

	// }

	err := netplan.WriteDefaultConfig()
	if err != nil {
		panic(err)
	}

	np, err := netplan.Load(netplanFile)
	if err != nil {
		panic(err)
	}
	fmt.Println(np)

}

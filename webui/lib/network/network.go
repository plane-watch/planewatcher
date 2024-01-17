package network

import (
	"net"
	"runtime"

	"github.com/vishvananda/netlink"
)

func init() {
	if runtime.GOOS != "linux" {
		panic("Unsupported on this operating system, as Linux netlink is required.")
	}
}

// isLoopback returns true if nf has net.FlagLoopback bit set
func isLoopback(nf net.Flags) bool {
	if nf&net.FlagLoopback == net.FlagLoopback { // bitwise AND to check for loopback
		return true
	} else {
		return false
	}
}

// returns a list of network devices (excludes loopback and bridge etc)
func GetNetworkDevices() ([]netlink.Link, error) {

	// get list if network interfaces
	ll, err := netlink.LinkList()
	if err != nil {
		return []netlink.Link{}, err
	}

	// remove loopback
	llDevicesOnly := []netlink.Link{}
	for _, l := range ll {
		if l.Type() == "device" {
			if !isLoopback(l.Attrs().Flags) {
				llDevicesOnly = append(llDevicesOnly, l)
			}
		}
	}

	// return sanitised list
	return llDevicesOnly, nil
}

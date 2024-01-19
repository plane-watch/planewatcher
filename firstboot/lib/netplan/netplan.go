package netplan

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/vishvananda/netlink"
	"gopkg.in/yaml.v2"
)

type Netplan struct {
	Network Network `yaml:"network"`
}

type Network struct {
	Version   int                 `yaml:"version"`
	Renderer  string              `yaml:"renderer,omitempty"`
	Ethernets map[string]Ethernet `yaml:"ethernets,omitempty"`
}

type Ethernet struct {
	Interface `yaml:",inline"`
}

type Interface struct {
	Addresses []string `yaml:"addresses,omitempty"`
	// DHCP4 defaults to true, so we must use a pointer to know if it was specified as false
	DHCP4       *bool       `yaml:"dhcp4,omitempty"`
	DHCP6       *bool       `yaml:"dhcp6,omitempty"`
	Gateway4    string      `yaml:"gateway4,omitempty"`
	Nameservers Nameservers `yaml:"nameservers,omitempty"`
	MTU         int         `yaml:"mtu,omitempty"`
	Routes      []Route     `yaml:"routes,omitempty"`
}

type Route struct {
	From   string `yaml:"from,omitempty"`
	OnLink *bool  `yaml:"on-link,omitempty"`
	Scope  string `yaml:"scope,omitempty"`
	Table  *int   `yaml:"table,omitempty"`
	To     string `yaml:"to,omitempty"`
	Type   string `yaml:"type,omitempty"`
	Via    string `yaml:"via,omitempty"`
	Metric *int   `yaml:"metric,omitempty"`
}

type Nameservers struct {
	Search    []string `yaml:"search,omitempty,flow"`
	Addresses []string `yaml:"addresses,omitempty,flow"`
}

func Load(filename string) (Netplan, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Netplan{}, err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return Netplan{}, err
	}

	fmt.Println(string(b))

	np := Netplan{}

	err = yaml.Unmarshal(b, &np)
	if err != nil {
		return Netplan{}, err
	}

	return np, err
}

// WriteDefaultConfig writes a default netplan yaml config with dchp4 enabled for all detected interfaces
func DefaultConfig(filename string) error {

	// open netlink file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// prep vars
	yes := true
	eths := make(map[string]Ethernet)

	// get ip link list
	ll, err := netlink.LinkList()
	if err != nil {
		return err
	}

	// for each link...
	for _, l := range ll {
		// if device (as opposed to bridge etc)
		if l.Type() == "device" {
			// if not loopback
			if !(l.Attrs().Flags&net.FlagLoopback == net.FlagLoopback) {
				// add interface
				eths[l.Attrs().Name] = Ethernet{
					Interface: Interface{
						DHCP4: &yes,
					},
				}
			}
		}
	}

	// prep netplan obj for marshalling to yaml
	np := Netplan{
		Network: Network{
			Version:   2,
			Renderer:  "networkd",
			Ethernets: eths,
		},
	}

	// marshall netplan obj to yaml
	out, err := yaml.Marshal(&np)
	if err != nil {
		return err
	}

	// write output
	_, err = f.Write(out)
	if err != nil {
		return err
	}

	// chmod
	err = os.Chmod(filename, 0600)
	if err != nil {
		return err
	}

	// apply
	err = ApplyImmediate()
	if err != nil {
		return err
	}

	return nil
}

// ApplyImmediate runs `netplan apply`
func ApplyImmediate() error {

	// prepare command
	c := exec.Command("netplan", "apply")

	// prepare stdout & stderr
	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}

	// start process
	err = c.Start()
	if err != nil {
		return err
	}

	// read stdout & stderr
	bStdout, err := io.ReadAll(stdout)
	if err != nil {
		return err
	}
	bStderr, err := io.ReadAll(stderr)
	if err != nil {
		return err
	}

	log := log.
		With().
		Str("stdout", string(bStdout)).
		Str("stderr", string(bStderr)).
		Logger()

	// wait for execution to finish
	err = c.Wait()
	if err != nil {
		log.
			Err(err).
			Msg("error running netplan apply")
		return err
	}

	log.
		Debug().
		Msg("ran netplan apply")

	return nil
}

// c := exec.Command("netplan", "try")
// stdin, err := c.StdinPipe()
// if err != nil {
// 	panic(err)
// }
// stdout, err := c.StdoutPipe()
// if err != nil {
// 	panic(err)
// }

// err = c.Start()
// if err != nil {
// 	panic(err)
// }
// time.Sleep(5 * time.Second)
// _, err = stdin.Write([]byte("\n"))
// if err != nil {
// 	panic(err)
// }
// b, err := io.ReadAll(stdout)
// if err != nil {
// 	panic(err)
// }
// err = c.Wait()
// if err != nil {
// 	panic(err)
// }

package webui

import (
	_ "embed"
	"firstboot/lib/netplan"
	"html/template"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

var (
	//go:embed netconfig.html
	tmplNetworkConfigHTML string

	// html template for network config
	tmplNetworkConfig *template.Template

	// netplan yaml file path
	netplanFile string
)

type networkConfig struct {
	Netplan   netplan.Netplan
	Interface map[string]netiface
}

type netiface struct {
	ipv4addr, ipv4mask, ipv4gw string
}

// defines the configuration for the Web UI service
type WebUI struct {
	ListenAddr  string
	NetplanFile string
}

func handleNetworkConfig(w http.ResponseWriter, r *http.Request) {
	var err error

	log := log.With().Str("netplan_yaml", netplanFile).Logger()

	// prep network config
	nc := networkConfig{}
	nc.Interface = make(map[string]netiface)

	// load netplan yaml
	nc.Netplan, err = netplan.Load(netplanFile)
	if err != nil {
		log.Err(err).Msg("error loading netplan yaml")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for iface := range nc.Netplan.Network.Ethernets {
		log := log.With().Str("iface", iface).Logger()
		// get "live" network config for each interface
		l, err := netlink.LinkByName(iface)
		if err != nil {
			log.Err(err).Msg("error getting interface information from netlink")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		addrs, err := netlink.AddrList(l, unix.AF_INET)
		if err != nil {
			log.Err(err).Msg("error getting interface addresses from netlink")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(addrs) < 1 {
			log.Err(err).Msg("no ipv4 addresses returned from netlink for interface")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(addrs) > 1 {
			log.Err(err).Msg("too many ipv4 addresses returned from netlink for interface")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		nc.Interface[iface] = netiface{
			ipv4addr: string(addrs[0].IP),
			ipv4mask: string(addrs[0].Mask),
		}

	}

	err = tmplNetworkConfig.Execute(w, nc)
	if err != nil {
		log.Err(err).Msg("error executing template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (conf *WebUI) Run() {
	var err error

	netplanFile = conf.NetplanFile

	// handle requests to network config page
	tmplNetworkConfig, err = template.New("NetworkConfig").Parse(tmplNetworkConfigHTML)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", handleNetworkConfig)

	err = http.ListenAndServe(conf.ListenAddr, nil)
	if err != nil {
		panic(err)
	}

}

package webui

import (
	_ "embed"
	"firstboot/lib/netplan"
	"html/template"
	"net/http"

	"github.com/rs/zerolog/log"
)

var (

	//go:embed netconfig.html
	tmplNetworkConfigHTML string

	tmplNetworkConfig *template.Template

	netplanFile string
)

type WebUIConfig struct {
	ListenAddr string

	NetplanFile string
}

func handleNetworkConfig(w http.ResponseWriter, r *http.Request) {
	log := log.With().Str("netplan_yaml", netplanFile).Logger()

	np, err := netplan.Load(netplanFile)
	if err != nil {
		log.Err(err).Msg("error loading netplan yaml")
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = tmplNetworkConfig.Execute(w, np)
	if err != nil {
		log.Err(err).Msg("error executing template")
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (conf *WebUIConfig) Run() {
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

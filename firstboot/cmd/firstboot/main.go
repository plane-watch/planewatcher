package main

import (
	"os"
	"time"

	"firstboot/lib/netplan"
	"firstboot/lib/webui"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	// set up logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.UnixDate})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// check if netplan file exists
	if !fileExists(netplanFile) {
		log.Debug().Str("netplan_config", netplanFile).Msg("generating firstrun config")
		err := netplan.DefaultConfig(netplanFile)
		if err != nil {
			panic(err)
		}
	}

	// // test
	// confirm := netplan.ApplyWithConfirmation(60)
	// time.Sleep(time.Second * 3)
	// err := confirm()
	// if err != nil {
	// 	panic(err)
	// }

	webUI := webui.WebUI{
		ListenAddr:  ":80",
		NetplanFile: netplanFile,
	}
	webUI.Run()

}

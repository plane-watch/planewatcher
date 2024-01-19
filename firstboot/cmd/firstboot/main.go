package main

import (
	"firstboot/lib/netplan"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
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
	if !fileExists(netplanFile) {
		err := netplan.WriteDefaultConfig(netplanFile)
		if err != nil {
			panic(err)
		}

		c := exec.Command("netplan", "try")
		stdin, err := c.StdinPipe()
		if err != nil {
			panic(err)
		}
		stdout, err := c.StdoutPipe()
		if err != nil {
			panic(err)
		}

		err = c.Start()
		if err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Second)
		_, err = stdin.Write([]byte("\n"))
		if err != nil {
			panic(err)
		}
		b, err := io.ReadAll(stdout)
		if err != nil {
			panic(err)
		}
		err = c.Wait()
		if err != nil {
			panic(err)
		}

		fmt.Println(string(b))
	}
}

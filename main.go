package main

import (
	"gioui.org/app"
	"github.com/partisiadev/partisiawallet/log"
	"github.com/partisiadev/partisiawallet/ui"
	"os"
	"os/exec"
	"strings"
	"time"
)

// FixTimezone https://github.com/golang/go/issues/20455
func FixTimezone() {
	out, err := exec.Command("/system/bin/getprop", "persist.sys.timezone").Output()
	if err != nil {
		return
	}
	z, err := time.LoadLocation(strings.TrimSpace(string(out)))
	if err != nil {
		return
	}
	time.Local = z
}

func main() {
	FixTimezone()
	go func() {
		if err := ui.Loop(); err != nil {
			log.Logger().Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

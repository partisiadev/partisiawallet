package main

import (
	"github.com/partisiadev/partisiawallet/app"
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
	app.Run()
}

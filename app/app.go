package app

import (
	"gioui.org/app"
	"github.com/partisiadev/partisiawallet/app/internal/state"
	"github.com/partisiadev/partisiawallet/app/internal/ui"
	log "github.com/sirupsen/logrus"
	"os"
)

func Run() {
	go func() {
		if err := state.Loop(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	go func() {
		if err := ui.Loop(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

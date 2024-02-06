package ui

import (
	"errors"
	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/ui/theme"
	"strconv"

	"github.com/partisiadev/partisiawallet/log"
	"time"
)

type FrameTiming struct {
	Start, End      time.Time
	FrameCount      int
	FramesPerSecond float64
}

var isRunning bool

func Loop() error {
	if isRunning {
		return errors.New("ui loop is already running")
	}
	isRunning = true
	w := app.NewWindow(app.Title("Multi Wallet"), app.Size(350, 600))
	uiManager := newAppManager(w)
	var ops op.Ops
	var deco widget.Decorations
	var decorated bool
	var title string
	option := app.StatusColor(theme.GlobalTheme.Theme().ContrastBg)
	w.Option(option)
	w.Option(app.Decorated(false))
	// backClickTag is meant for tracking db's backClick action, specially on mobile
	var backClickTag struct{}
	timingWindow := time.Second
	var timings []FrameTiming
	frameCounter := 0
	timingStart := time.Time{}
	for {
		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			log.Logger().Errorln("system.DestroyEvent called", e.Err)
			return e.Err
		case app.ConfigEvent:
			decorated = e.Config.Decorated
			title = e.Config.Title
		case system.FrameEvent:
			uiManager.insets = e.Insets
			e.Insets = system.Insets{}
			gtx := layout.NewContext(&ops, e)
			op.InvalidateOp{}.Add(gtx.Ops)
			if timingStart == (time.Time{}) {
				timingStart = gtx.Now
			}
			if interval := gtx.Now.Sub(timingStart); interval >= timingWindow {
				timings = append(timings, FrameTiming{
					Start:           timingStart,
					End:             gtx.Now,
					FrameCount:      frameCounter,
					FramesPerSecond: float64(frameCounter) / interval.Seconds(),
				})
				frameCounter = 0
				timingStart = gtx.Now
			}
			for _, event := range gtx.Events(&backClickTag) {
				switch e := event.(type) {
				case key.Event:
					switch e.Name {
					case key.NameBack:
						uiManager.Router().PopUp()
					}
				}
			}
			// Listen to back command only when uiManager.pagesStack is greater than 1,
			//  so we can pop up page else we want the android's default behavior
			if uiManager.Router().StackSize() > 1 {
				key.InputOp{Tag: &backClickTag, Keys: key.NameBack}.Add(gtx.Ops)
			}
			uiManager.metric = gtx.Metric
			uiManager.constraints = gtx.Constraints
			uiManager.Layout(gtx)
			if !decorated {
				w.Perform(deco.Update(gtx))
				uiManager.decoratedSize = material.Decorations(theme.GlobalTheme.Theme(), &deco, ^system.Action(0), title).Layout(gtx)
			} else {
				uiManager.decoratedSize = layout.Dimensions{}
			}
			e.Frame(gtx.Ops)
			for _, timing := range timings {
				_ = timing
				txt2 := strconv.FormatFloat(timing.FramesPerSecond, 'f', 2, 64)
				log.Logger().Println(txt2)
			}
			frameCounter++
		case system.StageEvent:
			if e.Stage == system.StagePaused {
				log.Logger().Infoln("window is running in background")
				uiManager.isStageRunning = false
			} else if e.Stage == system.StageRunning {
				log.Logger().Infoln("window is running in foreground")
				uiManager.isStageRunning = true
			}
		default:
		}
	}
}

package view

import (
	"fmt"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"image"
	"math"
	"time"
)

const minimumSwipeSpeed = float64(0.75)

type Swiper struct {
	dragging      bool
	dragStart     float32
	dragOffset    float32
	selectedIndex int
	layout.Axis

	thresholdTime time.Time //
	AllowInfSwipe bool
	// pixel per second
	swipeSpeed float64
}

func (s *Swiper) Layout(gtx layout.Context, length int, element layout.ListElement) layout.Dimensions {
	rec := op.Record(gtx.Ops)
	currDim := element(gtx, s.selectedIndex)
	rec.Stop()
	viewDim := gtx.Constraints.Max
	maxDim := viewDim
	if s.Axis.Convert(currDim.Size).X > s.Axis.Convert(viewDim).X {
		maxDim = currDim.Size
	}
	areaStack := clip.Rect(image.Rectangle{Max: maxDim}).Push(gtx.Ops)
	s.handleDrag(gtx, length)
	maxDrag := float32(s.Axis.Convert(maxDim).X)
	if s.dragging {
		if s.dragOffset >= maxDrag {
			s.dragOffset = maxDrag
		} else if s.dragOffset <= -maxDrag {
			s.dragOffset = -maxDrag
		}
	}

	if !s.dragging && s.dragOffset != 0 && maxDrag != 0 {
		var delta time.Duration
		if !s.thresholdTime.IsZero() {
			now := gtx.Now
			delta = now.Sub(s.thresholdTime)
			s.thresholdTime = now
		}
		movement := float32(math.Abs(float64(float32(s.swipeSpeed) * float32(delta.Milliseconds()))))
		if s.dragOffset < 0 {
			if s.swipeSpeed > minimumSwipeSpeed {
				s.dragOffset -= movement
			} else {
				s.dragOffset += movement * float32(minimumSwipeSpeed) * float32(delta.Milliseconds())
			}
			if s.dragOffset <= -maxDrag {
				s.selectedIndex++
				s.dragOffset = 0
			}
			if s.dragOffset >= 0 {
				s.dragOffset = 0
			}
		} else if s.dragOffset > 0 {
			if s.swipeSpeed > minimumSwipeSpeed {
				s.dragOffset += movement
			} else {
				s.dragOffset -= movement * float32(minimumSwipeSpeed) * float32(delta.Milliseconds())
			}
			if s.dragOffset >= maxDrag {
				s.selectedIndex--
				s.dragOffset = 0
			}
			if s.dragOffset <= 0 {
				s.dragOffset = 0
			}
		}

		if s.selectedIndex < 0 {
			s.selectedIndex = length - 1
		} else if s.selectedIndex > length-1 {
			s.selectedIndex = 0
		}
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	defer areaStack.Pop()
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			switch s.Axis {
			case layout.Vertical:
				defer op.Offset(image.Point{Y: int(s.dragOffset)}).Push(gtx.Ops).Pop()
			case layout.Horizontal:
				fallthrough
			default:
				defer op.Offset(image.Point{X: int(s.dragOffset)}).Push(gtx.Ops).Pop()
			}
			gtx.Constraints.Min = maxDim
			return element(gtx, s.selectedIndex)
		}),
		// Left Widget
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			switch s.Axis {
			case layout.Vertical:
				defer op.Offset(image.Point{Y: int(s.dragOffset - float32(maxDim.Y))}).Push(gtx.Ops).Pop()
			case layout.Horizontal:
				fallthrough
			default:
				defer op.Offset(image.Point{X: int(s.dragOffset - float32(maxDim.X))}).Push(gtx.Ops).Pop()
			}
			gtx.Constraints.Min = maxDim
			selectedIndex := s.selectedIndex - 1
			if selectedIndex < 0 {
				selectedIndex = length - 1
			} else if selectedIndex > length-1 {
				selectedIndex = 0
			}
			return element(gtx, selectedIndex)
		}),
		// Right Widget
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			switch s.Axis {
			case layout.Vertical:
				defer op.Offset(image.Point{Y: int(s.dragOffset + float32(maxDim.Y))}).Push(gtx.Ops).Pop()
			case layout.Horizontal:
				fallthrough
			default:
				defer op.Offset(image.Point{X: int(s.dragOffset + float32(maxDim.X))}).Push(gtx.Ops).Pop()
			}
			gtx.Constraints.Min = maxDim
			selectedIndex := s.selectedIndex + 1
			if selectedIndex < 0 {
				selectedIndex = length - 1
			} else if selectedIndex > length-1 {
				selectedIndex = 0
			}
			return element(gtx, selectedIndex)
		}),
	)
}

func (s *Swiper) handleDrag(gtx layout.Context, length int) {
	pointer.InputOp{
		Tag:   s,
		Kinds: pointer.Press | pointer.Release | pointer.Cancel | pointer.Drag,
	}.Add(gtx.Ops)
	for _, e := range gtx.Events(s) {
		switch e := e.(type) {
		case pointer.Event:
			posLength := e.Position.X
			if s.Axis == layout.Vertical {
				posLength = e.Position.Y
			}
			switch e.Kind {
			case pointer.Press:
				s.dragStart = posLength
				s.dragOffset = 0
				s.dragging = true
				s.thresholdTime = gtx.Now
			case pointer.Scroll:
				fmt.Printf("%#v\n", e)
			case pointer.Drag:
				//fmt.Printf("%#v\n", e)
				s.dragOffset = posLength - s.dragStart
				s.dragging = true
				if s.selectedIndex == 0 && !s.AllowInfSwipe {
					if s.dragOffset > 0 {
						s.dragOffset = 0
					}
				}
				if s.selectedIndex == length-1 && !s.AllowInfSwipe {
					if s.dragOffset < 0 {
						s.dragOffset = 0
					}
				}
			case pointer.Release, pointer.Cancel:
				if s.dragOffset != 0 {
					s.swipeSpeed = math.Abs(float64(s.dragOffset) / float64(gtx.Now.Sub(s.thresholdTime).Milliseconds()))
				}
				s.thresholdTime = gtx.Now
				s.dragging = false
			default:
				panic("unhandled default case")
			}
		}
	}
}

package main

import (
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit" // (dp,sp,px) реализует независимые от устройства единицы и значения? dp
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var versionProg string

type C = layout.Context
type D = layout.Dimensions

// root level, outside main ()
var progress float32
var progressIncrementer chan float32

func main() {
	log.Println("Версия: ", versionProg)
	// Setup a separate channel to provide ticks to increment progress
	progressIncrementer = make(chan float32)
	go func() {
		for {
			time.Sleep(time.Second / 25)
			progressIncrementer <- 0.004
		}
	}()
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Написано на GO"),
			app.Size(unit.Dp(400), unit.Dp(600)), //  x-y
			// MaxSize MinSize Both
		)

		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
func draw(w *app.Window) error {
	// ...
	// блок переменных не должен быть в Events !
	// ops are the operations from the UI
	var ops op.Ops

	// startButton is a clickable widget
	var startButton widget.Clickable

	// th defnes the material design style
	th := material.NewTheme(gofont.Collection())
	// is the egg boiling?
	var boiling bool
	// listen for events in the window.
	for e := range w.Events() {

		// detect what type of event
		switch e := e.(type) {

		// this is sent when the application should re-render.
		case system.FrameEvent:

			// ...
			gtx := layout.NewContext(&ops, e)
			// Let's try out the flexbox layout concept
			if startButton.Clicked() {
				boiling = !boiling
			}
			// Let's try out the flexbox layout concept
			layout.Flex{
				// Vertical alignment, from top to bottom
				Axis: layout.Vertical,
				// Empty space is left at the start, i.e. at the top
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx C) D {
						bar := material.ProgressBar(th, progress) // Here progress is used
						return bar.Layout(gtx)
					},
				),
				layout.Rigid(
					func(gtx C) D {
						// ONE: First define margins around the button using layout.Inset ...
						margins := layout.Inset{
							Top:    unit.Dp(25),
							Bottom: unit.Dp(25),
							Right:  unit.Dp(35),
							Left:   unit.Dp(35),
						}
						// TWO: ... then we lay out those margins ...
						return margins.Layout(gtx,
							// THREE: ... and finally within the margins, we ddefine and lay out the button
							func(gtx C) D {
								var text string
								if !boiling {
									text = "Start"
								} else {
									text = "Stop"
								}
								btn := material.Button(th, &startButton, text)
								return btn.Layout(gtx)
							},
						)
					},
				),
			)
			e.Frame(gtx.Ops)
			// this is sent when the application is closed.
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}
func draw1(w *app.Window) error {
	// ops are the operations from the UI
	var ops op.Ops

	// startButton is a clickable widget
	var startButton widget.Clickable

	// th defnes the material design style
	th := material.NewTheme(gofont.Collection())
	// listen for events in the window.
	for e := range w.Events() {
		// detect what type of event
		switch e := e.(type) {

		// this is sent when the application should re-render.
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			// btn := material.Button(th, &startButton, "Start")
			// btn.Layout(gtx)
			//-------
			// Let's try out the flexbox layout concept:
			layout.Flex{
				// Vertical alignment, from top to bottom
				Axis: layout.Vertical,
				// Empty space is left at the start, i.e. at the top
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				// We insert two rigid elements:
				// First a button ...
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &startButton, "Start")
						return btn.Layout(gtx)
					},
				),
				// ... then an empty spacer
				layout.Rigid(
					// The height of the spacer is 25 Device independent pixels
					layout.Spacer{Height: unit.Dp(25)}.Layout,
				),
			)
			//--------
			e.Frame(gtx.Ops)
		}
	}
	return nil
}

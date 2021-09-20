package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit" // (dp,sp,px) реализует независимые от устройства единицы и значения? dp
	"github.com/ark-go/giostart/internal"
)

var versionProg string

type C = layout.Context
type D = layout.Dimensions

// root level, outside main ()
// var progress float32
// var progressIncrementer chan bool

func main() {
	log.Println("Версия: ", versionProg)
	// Setup a separate channel to provide ticks to increment progress
	// progressIncrementer = make(chan bool)
	// go func() {
	// 	for {
	// 		time.Sleep(time.Second / 25) // уменьшает нагрузку, меньше перерисовка экрана
	// 		progressIncrementer <- true
	// 	}
	// }()
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Написано на GO"),
			app.Size(unit.Dp(400), unit.Dp(600)), //  x-y
			app.MinSize(unit.Dp(200), unit.Dp(300)),
			// MaxSize MinSize Both
		)

		if err := internal.Draw3(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

// func layoutWidget(ctx layout.Context, width, height int) layout.Dimensions {
// 	return layout.Dimensions{
// 		Size: image.Point{
// 			X: width,
// 			Y: height,
// 		},
// 	}
// }

// func drawOld(w *app.Window) error {
// 	// ops are the operations from the UI
// 	var ops op.Ops

// 	// startButton is a clickable widget
// 	var startButton widget.Clickable

// 	// th defnes the material design style
// 	th := material.NewTheme(gofont.Collection())
// 	// listen for events in the window.
// 	for e := range w.Events() {
// 		// detect what type of event
// 		switch e := e.(type) {

// 		// this is sent when the application should re-render.
// 		case system.FrameEvent:
// 			gtx := layout.NewContext(&ops, e)
// 			// btn := material.Button(th, &startButton, "Start")
// 			// btn.Layout(gtx)
// 			//-------
// 			// Let's try out the flexbox layout concept:
// 			layout.Flex{
// 				// Vertical alignment, from top to bottom
// 				Axis: layout.Vertical,
// 				// Empty space is left at the start, i.e. at the top
// 				Spacing: layout.SpaceStart,
// 			}.Layout(gtx,
// 				// We insert two rigid elements:
// 				// First a button ...
// 				layout.Rigid(
// 					func(gtx layout.Context) layout.Dimensions {
// 						btn := material.Button(th, &startButton, "Start")
// 						return btn.Layout(gtx)
// 					},
// 				),
// 				// ... then an empty spacer
// 				layout.Rigid(
// 					// The height of the spacer is 25 Device independent pixels
// 					layout.Spacer{Height: unit.Dp(25)}.Layout,
// 				),
// 			)
// 			//--------
// 			e.Frame(gtx.Ops)
// 		}
// 	}
// 	return nil
// }

package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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
			time.Sleep(time.Second / 25) // уменьшает нагрузку, меньше перерисовка экрана
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
	for {
		select {
		// listen for events in the window.
		case e := <-w.Events():

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
						//! круг
						// func(gtx C) D {
						// 	circle := clip.Circle{
						// 		// Hard coding the x coordinate. Try resizing the window
						// 		Center: f32.Point{X: 200, Y: 200},
						// 		// Soft coding the x coordinate. Try resizing the window
						// 		// Center: f32.Point{X: float32(gtx.Constraints.Max.X) / 2, Y: 200},
						// 		Radius: 120,
						// 	}.Op(gtx.Ops)
						// 	color := color.NRGBA{R: 200, A: 255}
						// 	paint.FillShape(gtx.Ops, color, circle)
						// 	d := image.Point{Y: 500}
						// 	return layout.Dimensions{Size: d}
						// },
						func(gtx C) D {
							// Draw a custom path, shaped like an egg
							var eggPath clip.Path
							op.Offset(f32.Pt(200, 150)).Add(gtx.Ops)
							eggPath.Begin(gtx.Ops)
							// Rotate from 0 to 360 degrees
							for deg := 0.0; deg <= 360; deg++ {

								// Egg math (really) at this brilliant site. Thanks!
								// https://observablehq.com/@toja/egg-curve
								// Convert degrees to radians
								rad := deg / 360 * 2 * math.Pi
								// Trig gives the distance in X and Y direction
								cosT := math.Cos(rad)
								sinT := math.Sin(rad)
								// Constants to define the eggshape
								a := 110.0
								b := 150.0
								d := 20.0
								// The x/y coordinates
								x := a * cosT
								y := -(math.Sqrt(b*b-d*d*cosT*cosT) + d*sinT) * sinT
								// Finally the point on the outline
								p := f32.Pt(float32(x), float32(y))
								// Draw the line to this point
								eggPath.LineTo(p)
							}
							// Close the path
							eggPath.Close()

							// Get hold of the actual clip
							eggArea := clip.Outline{Path: eggPath.End()}.Op()

							// Fill the shape
							// color := color.NRGBA{R: 255, G: 239, B: 174, A: 255}
							color := color.NRGBA{R: 255, G: uint8(239 * (1 - progress)), B: uint8(174 * (1 - progress)), A: 255}
							paint.FillShape(gtx.Ops, color, eggArea)

							d := image.Point{Y: 375}
							return layout.Dimensions{Size: d}
						},
					),
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
										text = "Старт"
									} else {
										text = "Стоп"
									}
									btn := material.Button(th, &startButton, text)
									btn.TextSize = unit.Dp(20)
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
		case p := <-progressIncrementer:
			if boiling && progress < 1 {
				progress += p
				w.Invalidate() // перерисовать экран
				//	op.InvalidateOp{At: gtx.Now.Add(time.Second / 25)}.Add(&ops)
			} else if progress >= 1 {
				progress = 0
				w.Invalidate()
			}
		}
	}

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

package internal

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type C = layout.Context
type D = layout.Dimensions

var progress float32
var progressIncrementer chan bool

//type cursor pointer.CursorName

func Draw1(w *app.Window) error {
	progressIncrementer = make(chan bool)
	go func() {
		for {
			time.Sleep(time.Second / 25) // уменьшает нагрузку, меньше перерисовка экрана
			progressIncrementer <- true
		}
	}()
	// ...
	// блок переменных не должен быть в Events !
	// ops are the operations from the UI
	var ops op.Ops

	// startButton is a clickable widget
	var startButton widget.Clickable
	// boilDurationInput is a textfield to input boil duration
	var boilDurationInput widget.Editor
	var boilDuration float32
	// is the egg boiling?
	var boiling bool
	// th defnes the material design style
	th := material.NewTheme(gofont.Collection())

	// listen for events in the window.
	//
	var textBox widget.Editor
	//var scrollBox widget.Scrollbar
	var tweetInput component.TextField
	var alg layout.Alignment
	var list widget.List
	type animation struct {
		start    time.Time
		duration time.Duration
	}

	list.AddDrag(&ops)
	//var list2 widget.Image
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

					// Resetting the boil
					//	if progress >= 1 {
					progress = 0
					//	}
					// Read from the input box
					//	if progress == 0 {
					inputString := boilDurationInput.Text()
					inputString = strings.TrimSpace(inputString)
					inputFloat, _ := strconv.ParseFloat(inputString, 32)
					boilDuration = float32(inputFloat)
					//	}
					// Start (or stop) the boil
					boiling = !boiling
				}
				// Let's try out the flexbox layout concept
				layout.Flex{
					// Vertical alignment, from top to bottom
					Axis:      layout.Vertical,
					Alignment: layout.End,
					// Empty space is left at the start, i.e. at the top
					//Spacing: layout.SpaceStart,
					//Spacing: layout,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						//	bar := material.ProgressBar(th, progress)

						return material.List(th, &list).Layout(gtx, 100, func(gtx C, index int) D {
							ed := material.Editor(th, &boilDurationInput, fmt.Sprintf(" - %d - ", index))
							return ed.Layout(gtx)
							// ed1 := ed.Layout(gtx)
							// ed1.Size = image.Point{
							// 	X: 80,
							// 	Y: 50,
							// }
							// return ed1

						})

					}),
					layout.Rigid(func(gtx C) D {
						margins := layout.Inset{
							Top:    unit.Dp(10),
							Right:  unit.Dp(10),
							Bottom: unit.Dp(10),
							Left:   unit.Dp(10),
						}
						if tweetInput.TextTooLong() {
							tweetInput.SetError("Too many characters")
						} else {
							tweetInput.ClearError()
						}
						tweetInput.CharLimit = 128
						tweetInput.Helper = "Твиты содержат ограниченное количество символов"
						tweetInput.Alignment = alg
						return margins.Layout(gtx,
							func(gtx C) D {
								return tweetInput.Layout(gtx, th, "Текст")
							},
						)

					}),
					layout.Rigid(
						func(gtx C) D {
							textBox.SingleLine = false
							textBox.Alignment = text.Start

							//textBox.SetText("Тест\nда")
							// Define insets ...
							margins := layout.Inset{
								Top:    unit.Dp(10),
								Right:  unit.Dp(10),
								Bottom: unit.Dp(10),
								Left:   unit.Dp(10),
							}
							margins2 := layout.Inset{
								Top:    unit.Dp(10),
								Right:  unit.Dp(10),
								Bottom: unit.Dp(10),
								Left:   unit.Dp(10),
							}

							// ... and borders ...
							border := widget.Border{
								Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
								CornerRadius: unit.Dp(3),
								Width:        unit.Dp(2),
							}
							// ... and material design ...
							ed := material.Editor(th, &textBox, "Введите текст")

							obj := margins.Layout(gtx, // в margin
								func(gtx C) D { // вставляем border
									//layoutWidget(gtx, 150, 150)
									return border.Layout(gtx,
										func(gtx C) D { // в border вставляем margin и элемент
											return margins2.Layout(gtx,
												func(gtx C) D {
													return ed.Layout(gtx)
												},
											)
										},
									)
								},
							)
							//obj.Size.Add(image.Point{Y: 150})
							return obj

							//	return layoutWidget(gtx, 50, 150)
							// ... before laying it out, one inside the other
							// return margins.Layout(gtx,
							// 	func(gtx C) D {
							// 		//	return border.Layout(gtx, ed.Layout)
							// 		return border.Layout(gtx, marginText)
							// 	},
							// )
						},
					),
					// layout.Rigid(
					// 	//! круг
					// 	func(gtx C) D {
					// 		circle := clip.Circle{
					// 			// Hard coding the x coordinate. Try resizing the window
					// 			Center: f32.Point{X: 200, Y: 200},
					// 			// Soft coding the x coordinate. Try resizing the window
					// 			// Center: f32.Point{X: float32(gtx.Constraints.Max.X) / 2, Y: 200},
					// 			Radius: 120,
					// 		}.Op(gtx.Ops)
					// 		color := color.NRGBA{R: 200, A: 255}
					// 		paint.FillShape(gtx.Ops, color, circle)
					// 		d := image.Point{Y: 500}
					// 		return layout.Dimensions{Size: d}
					// 	},
					// ),
					// The inputbox
					layout.Rigid(
						func(gtx C) D {
							// Define characteristics of the input box
							//	boilDurationInput.SingleLine = true
							boilDurationInput.Alignment = text.Middle

							// Count down the text when boiling
							if boiling && progress < 1 {
								boilRemain := (1 - progress) * boilDuration
								// Format to 1 decimal.
								//Использование старого доброго трюка умножения на 10 и деления на 10 для получения округленных значений с 1 десятичным знаком
								inputStr := fmt.Sprintf("%.1f", math.Round(float64(boilRemain)*10)/10)
								boilDurationInput.SetText(inputStr)
							}

							// Define insets ...
							margins := layout.Inset{
								Top:    unit.Dp(0),
								Right:  unit.Dp(170),
								Bottom: unit.Dp(40),
								Left:   unit.Dp(170),
							}
							// ... and borders ...
							border := widget.Border{
								Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
								CornerRadius: unit.Dp(3),
								Width:        unit.Dp(2),
							}
							// ... and material design ...
							ed := material.Editor(th, &boilDurationInput, "sec")
							// ... before laying it out, one inside the other
							return margins.Layout(gtx,
								func(gtx C) D {
									return border.Layout(gtx, ed.Layout)
								},
							)
						},
					),
					layout.Rigid(
						func(gtx C) D {
							bar := material.ProgressBar(th, progress) // Here progress is used
							if boiling && progress < 1 {
								op.InvalidateOp{At: gtx.Now.Add(time.Second / 100)}.Add(&ops)
							} else if progress >= 1 {
								progress = 0
								op.InvalidateOp{At: gtx.Now.Add(time.Second / 100)}.Add(&ops)
							}
							return bar.Layout(gtx)
						},
					),
					layout.Flexed(1, layout.Widget( // заполняет все оставшеееся пространство
						func(gtx C) D {
							cp := material.Caption(th, "пусто")
							return cp.Layout(gtx)
							// border := widget.Border{
							// 	Color:        color.NRGBA{R: 255, G: 204, B: 204, A: 255},
							// 	CornerRadius: unit.Dp(3),
							// 	Width:        unit.Dp(2),
							// }
							// return border.Layout(gtx,
							// 	func(gtx C) D {
							// 		cp := material.Caption(th, "")
							// 		return cp.Layout(gtx)
							// 		// return layout.Dimensions{
							// 		// 	Size: image.Point{
							// 		// 		// X: 115,
							// 		// 		// Y: 115,
							// 		// 	},
							// 		// }
							// 	})
						},
					)),
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
									btn.TextSize = unit.Dp(30)
									btn.CornerRadius = unit.Dp(20)

									dims := btn.Layout(gtx) // размеры?
									//! так установим курсор для кнопки , т.е. между dims
									//dims := material.Button(th, &startButton, text).Layout(gtx)
									pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
									return dims
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
		case <-progressIncrementer:
			if boiling && progress < 1 {
				//progress += p
				progress += 1.0 / 25.0 / boilDuration
				//op.InvalidateOp{}.Add(&ops) // Experimented, but couldn't make it perform
				//w.Invalidate() // перерисовать экран
				//	op.InvalidateOp{At: gtx.Now.Add(time.Second / 25)}.Add(&ops)
				// } else if progress >= 1 {
				// 	progress = 0
				// 	//op.InvalidateOp{}.Add(&ops)
				// 	//w.Invalidate()
			}
		}
	}

}

package internal

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

//type cursor pointer.CursorName

func Draw3(w *app.Window) error {
	type C = layout.Context
	type D = layout.Dimensions

	var progress float32
	progressIncrementer := make(chan bool)
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
	lineEditor := &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	editor := new(widget.Editor)
	// listen for events in the window.
	//
	//var scrollBox widget.Scrollbar
	var tweetInput component.TextField
	var alg layout.Alignment
	var list widget.List
	list1 := &widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
	list.AddDrag(&ops)
	xd := 40
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

				if startButton.Clicked() {

					progress = 0

					inputString := boilDurationInput.Text()
					inputString = strings.TrimSpace(inputString)
					inputFloat, _ := strconv.ParseFloat(inputString, 32)
					boilDuration = float32(inputFloat)
					//	}
					// Start (or stop) the boil
					boiling = !boiling
				}
				// Let's try out the flexbox layout concept
				// layout.Flex{
				// 	// Vertical alignment, from top to bottom
				// 	Axis:      layout.Vertical,
				// 	Alignment: layout.End,
				// 	// Empty space is left at the start, i.e. at the top
				// 	//Spacing: layout.SpaceStart,
				// 	//Spacing: layout,
				// }.Layout(gtx,
				widgets := []layout.Widget{
					material.H5(th, "ГО рисуем").Layout,
					func(gtx C) D {
						gtx.Constraints.Max.Y = gtx.Px(unit.Dp(100))
						return material.Editor(th, editor, "Hint").Layout(gtx)
					},
					//	layout.SpaceSides,
					//
					func(gtx C) D {
						e := material.Editor(th, lineEditor, "Hint")
						e.Font.Style = text.Italic
						border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(2)}
						return border.Layout(gtx, func(gtx C) D {
							return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
						})
					},
					//! круг
					func(gtx C) D {
						radius := float32(40)
						if xd > gtx.Constraints.Max.X-40 {
							xd = 40
							time.Sleep(time.Microsecond * 100)
						}
						xd++
						op.InvalidateOp{}.Add(gtx.Ops)
						log.Println(xd)
						circle := clip.Circle{
							// Hard coding the x coordinate. Try resizing the window
							//Center: f32.Point{X: 40, Y: 40},
							// Soft coding the x coordinate. Try resizing the window
							//Center: f32.Point{X: float32(gtx.Constraints.Max.X) / 2, Y: 40}, // по центру
							Center: f32.Point{X: float32(xd), Y: 40}, // по центру
							Radius: radius,
						}.Op(gtx.Ops)
						color := color.NRGBA{R: 200, A: 255}
						paint.FillShape(gtx.Ops, color, circle)

						d := image.Point{Y: int(radius * 2)}
						return layout.Dimensions{Size: d}
					},
					func(gtx C) D {
						if tweetInput.TextTooLong() {
							tweetInput.SetError("Too many characters")
						} else {
							tweetInput.ClearError()
						}
						tweetInput.CharLimit = 128
						tweetInput.Helper = "Твиты содержат ограниченное количество символов"
						tweetInput.Alignment = alg
						// return layout.UniformInset(unit.Dp(18)).Layout(gtx, func(gtx C) D {
						// 	return tweetInput.Layout(gtx, th, "Текст")
						// })

						return tweetInput.Layout(gtx, th, "Текст")

					},
					// The inputbox

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
							Bottom: unit.Dp(10),
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
					//func(gtx C) D {
					//dims :=
					// layout.Flexed(1, layout.Widget( // заполняет все оставшеееся пространство
					// 	func(gtx C) D {
					// 		cp := material.Caption(th, "пусто")
					// 		return cp.Layout(gtx)
					// 	})),
					//return dims
					//},
					// layout.Flexed(1, layout.Widget( // заполняет все оставшеееся пространство
					// 	func(gtx C) D {
					// 		cp := material.Caption(th, "пусто")
					// 		return cp.Layout(gtx)
					// 		// border := widget.Border{
					// 		// 	Color:        color.NRGBA{R: 255, G: 204, B: 204, A: 255},
					// 		// 	CornerRadius: unit.Dp(3),
					// 		// 	Width:        unit.Dp(2),
					// 		// }
					// 		// return border.Layout(gtx,
					// 		// 	func(gtx C) D {
					// 		// 		cp := material.Caption(th, "")
					// 		// 		return cp.Layout(gtx)
					// 		// 		// return layout.Dimensions{
					// 		// 		// 	Size: image.Point{
					// 		// 		// 		// X: 115,
					// 		// 		// 		// Y: 115,
					// 		// 		// 	},
					// 		// 		// }
					// 		// 	})
					// 	},
					// )),
				}
				//)
				// так мы все элементы вставляем в List со скролом
				// material.List(th, list1).Layout(gtx, len(widgets), func(gtx C, i int) D {
				// 	return layout.UniformInset(unit.Dp(10)).Layout(gtx, widgets[i]) // тут margin 10???
				// })

				layout.Flex{
					// 	// Vertical alignment, from top to bottom
					Axis:      layout.Vertical,
					Alignment: layout.End,
					// Empty space is left at the start, i.e. at the top
					//Spacing: layout.SpaceStart,
					//Spacing: layout,
				}.Layout(gtx,
					// Будет сжиматься наш list в оставшемся свободном пространстве
					layout.Flexed(1, layout.Widget( // заполняет все оставшеееся пространство
						func(gtx C) D {
							// cp := material.Caption(th, "пусто")
							// return cp.Layout(gtx)
							return material.List(th, list1).Layout(gtx, len(widgets), func(gtx C, i int) D {
								return layout.UniformInset(unit.Dp(10)).Layout(gtx, widgets[i]) // UniformInset общий margin 10???
							})
						})),
					// обычный flex по размеру компонента
					layout.Rigid(func(gtx C) D {
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

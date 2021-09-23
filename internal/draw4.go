package internal

import (
	"image"
	"image/color"
	"log"
	"runtime"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

//type cursor pointer.CursorName
type C = layout.Context
type D = layout.Dimensions
type page struct {
	list    *widget.List
	list2   *widget.List
	ops     *op.Ops
	gtx     layout.Context
	editor  *widget.Editor
	th      *material.Theme
	colOn   bool
	scrollY float32
	icon    *widget.Icon
	split   *Split
}

func Draw4(w *app.Window) error {
	// иконка для кнопки
	ic, err := widget.NewIcon(icons.CommunicationCall)
	if err != nil {
		log.Fatal(err)
	}
	icon = ic
	// ^^^^^^ иконка для кнопки  ^^^^^^^

	// y-position for scroll text
	// y-position for red highlight bar
	var highlightY float32 = 78
	Page := page{}
	Page.icon = icon
	Page.list = &widget.List{
		//Scrollbar: widget.Scrollbar{},
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
	Page.list2 = &widget.List{
		//Scrollbar: widget.Scrollbar{},
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
	Page.ops = &op.Ops{}
	Page.editor = new(widget.Editor)
	Page.th = material.NewTheme(gofont.Collection())
	Page.editor.SetText(texttest)
	Page.split = &Split{}
	log.Println("Процессоров:", runtime.NumCPU())
	log.Println("Проц:", MaxParallelism())
	for {
		select {
		// Слушать события в окне.
		case e := <-w.Events():
			// Определить, какой тип события
			switch e := e.(type) {
			// событие мыши
			case pointer.Event:
				// мышь колесико
				if e.Type == pointer.Scroll { // колесико
					stepSize := e.Scroll.Y // вниз в плюс, вверх в минус (+/-120)
					Page.scrollY = Page.scrollY - stepSize
					highlightY = highlightY + stepSize
					log.Println(stepSize, "|", highlightY)
					w.Invalidate()
				}
				// мышь кнопки
				if e.Type == pointer.Press {
					if e.Buttons == pointer.ButtonSecondary { // правая левая средняя
						log.Println("Правая", e.Position.X, e.Position.Y)
						log.Println("gtx", Page.gtx.Constraints.Min.X, Page.gtx.Constraints.Min.Y, Page.gtx.Constraints.Max.X, Page.gtx.Constraints.Max.Y)
					}
				}
			// Это отправляется, когда приложение должно повторно выполнить рендеринг.
			case system.FrameEvent:
				Page.gtx = layout.NewContext(Page.ops, e)
				//log.Println("start", Page.gtx)
				Page.render()
				draw1a(Page.gtx.Ops)
				draw1b(Page.gtx.Ops)

				e.Frame(Page.gtx.Ops) // Отрисовываем, без этого не рисует
			case system.DestroyEvent:
				return e.Err
				// событие клавиши
			case key.Event:
				if e.State == key.Press {
					if e.Name == "U" || e.Name == "K" {
						Page.colOn = !Page.colOn
						log.Println("Нажата U")
					}
					//op.InvalidateOp{}.Add(Page.gtx.Ops)
					w.Invalidate()
				}
				//	default:
				//		log.Println("gtx", Page.gtx.Constraints.Min.X, Page.gtx.Constraints.Max.X,Page.gtx.Constraints.Min.Y,Page.gtx.Constraints.Max.Y)
			}
		}

	}
}

func (p *page) render() {
	// Background
	if p.colOn {
		//красим всё разом
		background := clip.Rect{
			Min: image.Pt(0, 0),
			Max: image.Pt(p.gtx.Constraints.Max.X, p.gtx.Constraints.Max.Y),
		}.Op()
		paint.FillShape(p.ops, color.NRGBA{R: 0xff, G: 0xfe, B: 0xe0, A: 0xff}, background)

	}

	//-----------
	widgets := []layout.Widget{
		// текстовое поле
		func(gtx C) D {

			// высота текстового поля
			gtx.Constraints.Max.Y = gtx.Px(unit.Dp(100))
			// цвет background текстового поля
			dms := material.Editor(p.th, p.editor, "Hint").Layout(gtx)
			background := clip.Rect{
				Min: image.Pt(0, 0),
				Max: image.Pt(dms.Size.X, dms.Size.Y),
			}.Op()
			// красим
			paint.FillShape(p.ops, color.NRGBA{R: 0x22, G: 0xff, B: 0xf0, A: 0xff}, background)
			return material.Editor(p.th, p.editor, "Hint").Layout(gtx)
		},
		// ну обычный текст Llabel
		func(gtx C) D {
			speech := material.Label(p.th, unit.Dp(float32(20)), "Текст label")
			speech.Alignment = 2 // Center
			return speech.Layout(gtx)
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(8))
			return in.Layout(gtx, iconAndTextButton{theme: p.th, icon: p.icon, word: "Кнопка", button: iconTextButton}.Layout)
		},
	}
	//----------------------------
	widgets3 := []layout.FlexChild{
		// текстовое поле
		layout.Rigid(func(gtx C) D {
			// высота текстового поля
			gtx.Constraints.Max.Y = gtx.Px(unit.Dp(100))
			// цвет background текстового поля
			dms := material.Editor(p.th, p.editor, "Hint").Layout(gtx)
			background := clip.Rect{
				Min: image.Pt(0, 0),
				Max: image.Pt(dms.Size.X, dms.Size.Y),
			}.Op()
			// красим
			paint.FillShape(p.ops, color.NRGBA{R: 0x22, G: 0xff, B: 0xf0, A: 0xff}, background)
			return material.Editor(p.th, p.editor, "Hint").Layout(gtx)
		}),
		// ну обычный текст Llabel
		layout.Rigid(func(gtx C) D {
			speech := material.Label(p.th, unit.Dp(float32(20)), "Текст label")
			speech.Alignment = 2 // Center
			return speech.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			in := layout.UniformInset(unit.Dp(8))
			return in.Layout(gtx, iconAndTextButton{theme: p.th, icon: p.icon, word: "Кнопка", button: iconTextButton}.Layout)
		}),
		layout.Rigid(func(gtx C) D {
			//	return layout.Dimensions{Size: image.Point{X: 100}}
			kk := material.Body1(p.th, "Привет Привет Привет Привет Привет ").Layout(gtx)
			return kk
		}),
	}
	_ = widgets
	_ = widgets3
	//-----------------------------
	// основное оно
	widgets2 := []layout.Widget{
		func(gtx C) D {

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,

				layout.Rigid(func(gtx C) D {

					return layout.Flex{
						Axis:      layout.Horizontal,
						Spacing:   layout.SpaceSides,
						Alignment: layout.Alignment(layout.Center),
					}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = 50
							rr := layout.Flex{
								Axis: layout.Vertical,
								//	Spacing: layout.SpaceSides,
							}.Layout(gtx, widgets3...)

							return rr
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Max.Y = gtx.Px(unit.Dp(150))
							gt := material.List(p.th, p.list2).Layout(gtx, len(widgets), func(gtx C, i int) D {
								return layout.UniformInset(unit.Dp(5)).Layout(gtx, widgets[i]) // UniformInset общий margin 5 у каждого элемента в списке???

							})
							//gtx.Constraints.Min.X = gtx.Constraints.Max.X + gt.Size.X
							return gt
							//return material.Body1(p.th, "7777777 7 7 7 77 7 7 7 7 7 7 \n99 777777\n5555555555 55555555555 5555555555555555 5555555555555 5555555555 555").Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {

					return layout.Flex{
						Axis:      layout.Horizontal,
						Spacing:   layout.SpaceSides,
						Alignment: layout.End,
					}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							pp := gtx
							pp.Constraints.Max.X = 100
							pp.Constraints.Max.Y = 150
							return material.H5(p.th, "ГО рисуем").Layout(pp)
						}),
						layout.Rigid(func(gtx C) D {
							pp := gtx
							pp.Constraints.Max.X = 60
							pp.Constraints.Max.Y = 150
							return material.H5(p.th, "ГО рисуем").Layout(pp)
						}),
					)
				}),

				layout.Rigid(func(gtx C) D {

					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Start,
					}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							// размеры этого флекса ?
							//	gtx = p.gtx //! ,,,?????
							//	gtx.Constraints.Max.Y = 200 //p.gtx.Constraints.Max.Y //! а можно и не ограничивать вероятно
							//! Max.Y ограничивается общим размером по умолчанию .т.е надо исправлять? хз
							//	gtx.Constraints.Min.Y = 10
							//layout.Center.Layout(gtx, material.H3(p.th, "Право1").Layout)

							p.split.LeftPanelWidth(35) // можно объявить за циклом при создании p.split
							dims := p.split.Layout(gtx,
								func(gtx C) D {
									return layout.N.Layout(gtx, material.H3(p.th, "Лево").Layout)
								}, func(gtx C) D {
									return layout.N.Layout(gtx, material.H3(p.th, "Право1").Layout)
								})
							return dims //D{Size: image.Pt(50, 100)}
						}),
					)
				}),
			)
		}}
	//------------------------------------------
	layout.Flex{ // основной флекс
		Axis:      layout.Vertical,
		Alignment: layout.End,
		// Пустое место остается в начале, т.е. вверху
		//Spacing: layout.SpaceStart,
		//Spacing: layout,
	}.Layout(p.gtx,
		layout.Flexed(1, layout.Widget( // флекс заполняет все оставшеееся пространство

			func(gtx C) D {

				// cp := material.Caption(th, "пусто")
				// return cp.Layout(gtx)
				return material.List(p.th, p.list).Layout(gtx, len(widgets2), func(gtx C, i int) D {
					//log.Println("list-1 start", gtx)
					dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, widgets2[i]) // UniformInset общий margin 10???
					//log.Println("list-1 end", dims)
					//dims.Size.Y = p.gtx.Constraints.Max.Y
					return dims
				})

			})),
	)

}

//-------------- --------------------------
const texttest string = `
Привет строка ----- 9 - 1
Привет строка - 2
Привет строка - 3
Привет строка - 4
Привет строка - 5`

type iconAndTextButton struct {
	theme  *material.Theme
	button *widget.Clickable
	icon   *widget.Icon
	word   string
}

// -------------------------------------- ----------  кнопка ----------------------
// кнопка с иконкой
var iconTextButton = new(widget.Clickable)
var icon *widget.Icon

func (b iconAndTextButton) Layout(gtx layout.Context) layout.Dimensions {
	dims := material.ButtonLayout(b.theme, b.button).Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
			iconAndLabel := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
			textIconSpacer := unit.Dp(5)

			layIcon := layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: textIconSpacer}.Layout(gtx, func(gtx C) D {
					var d D
					if b.icon != nil {
						size := gtx.Px(unit.Dp(56)) - 2*gtx.Px(unit.Dp(16))
						gtx.Constraints = layout.Exact(image.Pt(size, size))
						d = b.icon.Layout(gtx, b.theme.ContrastFg)
					}
					return d
				})
			})

			layLabel := layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx C) D {
					l := material.Body1(b.theme, b.word)
					l.Color = b.theme.Palette.ContrastFg
					return l.Layout(gtx)
				})
			})

			return iconAndLabel.Layout(gtx, layIcon, layLabel)
		})

	})
	//		dims := btn.Layout(gtx) // размеры?
	//! так установим курсор для кнопки , т.е. между dims
	pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
	return dims
}

func draw1a(ops *op.Ops) {
	defer op.Save(ops).Load()
	clip.Rect{Min: image.Pt(20, 30), Max: image.Pt(100, 100)}.Add(ops)
	paint.ColorOp{Color: color.NRGBA{R: 0x80, A: 0xFF}}.Add(ops)
	paint.PaintOp{}.Add(ops)
}
func draw1b(ops *op.Ops) {
	defer op.Save(ops).Load()
	clip.Rect{Min: image.Pt(60, 60), Max: image.Pt(150, 180)}.Add(ops)
	paint.ColorOp{Color: color.NRGBA{G: 0x80, A: 0xFF}}.Add(ops)
	paint.PaintOp{}.Add(ops)
}
func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

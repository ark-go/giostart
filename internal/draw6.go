package internal

import (
	"log"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type pageMain struct {
	ops      *op.Ops
	gtx      layout.Context
	th       *material.Theme
	list     *widget.List
	split    *Split
	mainList []layout.Widget
	colOn    bool
	on       bool
}

var Page *pageMain

func init() {
	// Page = pageMain{}
	// Page.getMainList()
}

func InitPage() {
	Page = &pageMain{
		ops: &op.Ops{},
		th:  material.NewTheme(gofont.Collection()),
		list: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		split: &Split{},
	}

	//	Page.getMainList()
}

func (p *pageMain) Draw(w *app.Window) error {
	p.getMainList()
	log.Println("Старт Draw")
	for {
		select {
		// Слушать события в окне.
		case e := <-w.Events():
			// Определить, какой тип события
			switch e := e.(type) {
			// событие мыши
			case pointer.Event:
				// Это отправляется, когда приложение должно повторно выполнить рендеринг.
			case system.FrameEvent:
				Page.gtx = layout.NewContext(Page.ops, e)
				Page.on = true
				Page.render()
				e.Frame(Page.gtx.Ops) // Отрисовываем, без этого не рисует
			case system.DestroyEvent:
				return e.Err
				// событие клавиши
			case key.Event:
				if e.State == key.Press {
					if e.Name == "U" {
						Page.colOn = !Page.colOn
						log.Println("Нажата U")
					}
					//op.InvalidateOp{}.Add(Page.gtx.Ops)
					w.Invalidate()
				}
				//	default:
			}
		}
	}
}

func (p *pageMain) getMainList() {
	log.Println("getMainList")
	mainWidget := []layout.Widget{
		func(gtx C) D {

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(
					func(gtx C) D {
						// return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						// 	layout.Rigid(func(gtx C) D {
						return layout.Flex{
							Axis: layout.Horizontal,
							//	Alignment: layout.Start,
							// Spacing:   layout.SpaceSides,
							// Alignment: layout.Alignment(layout.Center),
						}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								// размеры этого флекса ?
								//gtx.Constraints.Max.Y = 200 //p.gtx.Constraints.Max.Y //! а можно и не ограничивать вероятно
								//! Max.Y ограничивается общим размером по умолчанию .т.е надо исправлять? хз
								//					gtx.Constraints.Min.Y = 170
								p.split.LeftPanelWidth(35) // можно объявить за циклом при создании p.split
								// cp := material.Caption(p.th, "\nпусто\111111111111111111111о\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm")
								// return cp.Layout(gtx)
								dims := p.split.Layout(gtx,
									func(gtx C) D {
										return layout.N.Layout(gtx, material.H3(p.th, "Лево").Layout)
									}, func(gtx C) D {
										return layout.N.Layout(gtx, material.H3(p.th, "Право1").Layout)
									})
								return dims //D{Size: image.Pt(50, 100)}
							}),
						)
					},
				),
			)
		},
	}

	p.mainList = mainWidget

}

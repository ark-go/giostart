package page

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
	"github.com/ark-go/giostart/internal/split"
)

//type cursor pointer.CursorName
type C = layout.Context
type D = layout.Dimensions
type TmainPage struct {
	MainList *[]layout.Widget
	List     *widget.List
	List2    *widget.List
	Ops      *op.Ops
	Gtx      layout.Context
	Editor   *widget.Editor
	Th       *material.Theme
	ColOn    bool
	ScrollY  float32
	Icon     *widget.Icon
	Split    *split.Split
	SplitV   *split.SplitV
	On       bool
}

//var MainPage *TmainPage

//func InitPage(p *TmainPage) *TmainPage {
func CreatePage() *TmainPage {
	return &TmainPage{
		Ops: &op.Ops{},
		Th:  material.NewTheme(gofont.Collection()),
		List: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		Split:  split.Create(40, 5), // &split.Split{},
		SplitV: split.CreateV(40, 5),
	}
}

func (p *TmainPage) Draw(w *app.Window) error {
	p.initMainList()
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
				p.Gtx = layout.NewContext(p.Ops, e)
				p.On = true
				p.Render()
				e.Frame(p.Gtx.Ops) // Отрисовываем, без этого не рисует
			case system.DestroyEvent:
				return e.Err
				// событие клавиши
			case key.Event:
				if e.State == key.Press {
					if e.Name == "U" {
						p.ColOn = !p.ColOn
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

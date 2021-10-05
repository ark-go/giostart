package internal

import ()

// type PageMain struct {
// 	ops      *op.Ops
// 	gtx      layout.Context
// 	th       *material.Theme
// 	list     *widget.List
// 	split    *Split
// 	mainList *[]layout.Widget
// 	colOn    bool
// 	on       bool
// }

//var Page *page.MainPage

func init() {
	//	 Page = PageMain{}

}

// func (p *PageMain) Draw(w *app.Window) error {
// 	p.initMainList()
// 	log.Println("Старт Draw")
// 	for {
// 		select {
// 		// Слушать события в окне.
// 		case e := <-w.Events():
// 			// Определить, какой тип события
// 			switch e := e.(type) {
// 			// событие мыши
// 			case pointer.Event:
// 				// Это отправляется, когда приложение должно повторно выполнить рендеринг.
// 			case system.FrameEvent:
// 				Page.gtx = layout.NewContext(Page.ops, e)
// 				Page.on = true
// 				Page.render()
// 				e.Frame(Page.gtx.Ops) // Отрисовываем, без этого не рисует
// 			case system.DestroyEvent:
// 				return e.Err
// 				// событие клавиши
// 			case key.Event:
// 				if e.State == key.Press {
// 					if e.Name == "U" {
// 						Page.colOn = !Page.colOn
// 						log.Println("Нажата U")
// 					}
// 					//op.InvalidateOp{}.Add(Page.gtx.Ops)
// 					w.Invalidate()
// 				}
// 				//	default:
// 			}
// 		}
// 	}
// }

package comp

import (
	"log"

	"gioui.org/layout"
	"github.com/ark-go/giostart/internal"
)

type C layout.Context
type D layout.Dimensions

func splitHorizontal(p *internal.PageMain, left layout.Widget, right layout.Widget) layout.FlexChild {

	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		// размеры этого флекса ?
		//gtx.Constraints.Max.Y = 200 //p.gtx.Constraints.Max.Y //! а можно и не ограничивать вероятно
		//! Max.Y ограничивается общим размером по умолчанию .т.е надо исправлять? хз
		//					gtx.Constraints.Min.Y = 170
		p.split.LeftPanelWidth(35) // можно объявить за циклом при создании p.split
		// cp := material.Caption(p.th, "\nпусто\111111111111111111111о\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm")
		// return cp.Layout(gtx)
		dims := p.split.Layout(gtx,
			func(gtx C) D {
				return left // layout.N.Layout(gtx, material.H3(p.th, "Лево").Layout)
			}, func(gtx C) D {
				return right //layout.N.Layout(gtx, material.H3(p.th, "Право1").Layout)
			})
		return dims //D{Size: image.Pt(50, 100)}
	})
}

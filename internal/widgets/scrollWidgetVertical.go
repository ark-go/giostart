package widgets

import (
	"image/color"
	"log"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

type TScrollWidget struct {
	Color color.NRGBA
}

var ScrollWidget *TScrollWidget

func init() {
	ScrollWidget = &TScrollWidget{}
}

var ScrollList *[]layout.FlexChild

func (sw *TScrollWidget) Layout(gtx layout.Context, children *[]layout.FlexChild) D {

	log.Println("TScroll count child:", len(*children), gtx.Constraints)

	dims := layout.Flex{ // вертикальный flex
		Axis: layout.Vertical,
	}.Layout(gtx, *children...) // детки
	log.Println("all size:", dims.Size)
	return dims
}

func init() {
	th := material.NewTheme(gofont.Collection())
	ScrollList = &[]layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return layout.N.Layout(gtx, material.H3(th, "Правый раз\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Center.Layout(gtx, material.H3(th, "Правый два\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Center.Layout(gtx, // это распологает виджет относительно всего доступного родительского виджета
				func(gtx C) D {
					return layout.N.Layout(gtx, material.H3(th, "Правый три\nВот: продолжение строки").Layout)
				})
		}),
	}
	_ = ScrollList
}

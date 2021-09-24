package internal

import (
	"log"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func (p *pageMain) render() {
	log.Println("Render", len(p.mainList), p.gtx)
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
				//gtx.Constraints.Min.X = 2500
				// cp := material.Caption(p.th, "пусто\nllllll")
				// return cp.Layout(gtx)
				return material.List(p.th, p.list).Layout(gtx, len(p.mainList), func(gtx C, i int) D {
					//log.Println("list-1 start", gtx)
					dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, p.mainList[i]) // UniformInset общий margin 10???
					//log.Println("list-1 end", dims)
					//dims.Size.Y = p.gtx.Constraints.Max.Y + 500
					log.Println("основной", dims)
					return dims
				})

			})),
	)
}

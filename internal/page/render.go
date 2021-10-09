package page

import (
	"image/color"
	"log"
	"runtime"
	"strconv"

	"gioui.org/layout"

	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/ark-go/giostart/internal/widgets"
)

var border *widgets.Border

func (p *TmainPage) Render() {
	log.Println("Render", p.Gtx.Constraints)
	//------------------------------------------
	layout.Flex{ // основной флекс
		Axis:      layout.Vertical,
		Alignment: layout.End,
		// Пустое место остается в начале, т.е. вверху
		//Spacing: layout.SpaceStart,
		//Spacing: layout,
	}.Layout(p.Gtx,
		layout.Flexed(1,
			func(gtx C) D {
				splitFlex.Init(40, 10)
				dims := splitFlex.Layout(gtx, leftList, rightList)
				return dims //D{Size: image.Pt(50, 100)}
			}),

		// layout.Flexed(1, layout.Widget( // флекс заполняет все оставшеееся пространство

		// 	func(gtx C) D {
		// 		//gtx.Constraints.Min.X = 2500
		// 		// cp := material.Caption(p.th, "пусто\nllllll")
		// 		// return cp.Layout(gtx)
		//return material.List(p.Th, p.List).Layout(gtx, len(*p.MainList), func(gtx C, i int) D {
		// 			dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, (*p.MainList)[i]) // UniformInset общий margin 10???
		// 			//dims.Size.Y = p.gtx.Constraints.Max.Y + 500
		// 			log.Println("основной", dims)
		// 			return dims
		// 		})

		// 	})),
		layout.Rigid(layout.Widget(
			func(gtx C) D {
				// поскольку в Border мы храним состояние (press key), чего-либо, то создать его надо только один раз
				// а раз нам нужен указатель, то можно проверить на nil
				//! или создавать его гдето еще?
				if border == nil {
					border = &widgets.Border{
						Color: color.NRGBA{R: 204, G: 204, B: 204, A: 255},
						Width: unit.Dp(2),
					}
				}
				// return border.Layout(gtx,
				// 	func(gtx C) D {
				// 		return layout.E.Layout(gtx, material.H3(p.Th, "Процессор: "+strconv.Itoa(runtime.NumCPU())).Layout)
				// 	})

				return layout.Center.Layout(gtx, // это распологает виджет относительно всего доступного родительского виджета
					func(gtx C) D {
						ggg := border.Layout(gtx, material.H6(p.Th, "Процессор: "+strconv.Itoa(runtime.NumCPU())).Layout)
						return ggg
					})
			})),
	)

}

// func setBorder(dim D, ops *op.Ops) D {
// 	defer op.Save(ops).Load()
// 	rec := f32.Rectangle{Min: f32.Point{X: 0, Y: 0}, Max: f32.Point{X: float32(dim.Size.X), Y: float32(dim.Size.Y)}}
// 	rr := clip.UniformRRect(rec, 0)
// 	rpath := rr.Path(ops)
// 	_ = rpath
// 	// str := clip.Stroke{
// 	// 	Path: rpath,
// 	// }
// 	// str.Op().Add(ops)
// 	// _ = color.NRGBA{R: 0x80, A: 0xFF}
// 	// //paint.ColorOp{Color: color.NRGBA{R: 0x80, A: 0xFF}}.Add(ops)
// 	// paint.PaintOp{}.Add(ops)
// 	// Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
// 	// CornerRadius: unit.Dp(3),
// 	// Width:        unit.Dp(2),
// 	clip.Rect{Min: image.Pt(60, 60), Max: image.Pt(150, 180)}.Add(ops)
// 	paint.ColorOp{Color: color.NRGBA{G: 0x80, A: 0xFF}}.Add(ops)
// 	paint.PaintOp{}.Add(ops)
// 	log.Println("Dim new", dim)
// 	return dim
// }

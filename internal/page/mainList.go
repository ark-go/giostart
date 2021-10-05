package page

import (
	"log"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

func (p *TmainPage) initMainList() {
	log.Println("initMainList")
	mainWidget := &[]layout.Widget{
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
								// gtx.Constraints.Min.Y = 170
								//.Split.LeftPanelWidth(35) // можно объявить за циклом при создании p.split
								// cp := material.Caption(p.th, "\nпусто\111111111111111111111о\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm\nпусто\nmmmmm")
								// return cp.Layout(gtx)+
								dims := p.Split.Layout(gtx,
									func(gtx C) D {
										return layout.N.Layout(gtx, material.H3(p.Th, "Лево\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
									}, func(gtx C) D {
										return layout.N.Layout(gtx, material.H3(p.Th, "Право\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
										//?  ///////////// 1вариант
										// gtx.Constraints.Max.Y = 600
										// //gtx.Constraints.Max = gtx.Constraints.Min
										// dims2 := p.SplitV.Layout(gtx,
										// 	func(gtx C) D {

										// 		log.Println("ttts:", gtx.Constraints)

										// 		return border.Layout(gtx, material.H3(p.Th, "Верх\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
										// 		//return layout.N.Layout(gtx, material.H3(p.Th, "Верх\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
										// 	}, func(gtx C) D {
										// 		return layout.N.Layout(gtx, material.H3(p.Th, "Низ\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)

										// 	})
										// return dims2
										//?  /////////// 2 Вариант
										//return layout.Rigid(func(gtx C) D {}
										//?  //////////
									})
								return dims //D{Size: image.Pt(50, 100)}
							}),
						)
					},
				),
			)
		},
	}

	p.MainList = mainWidget

}

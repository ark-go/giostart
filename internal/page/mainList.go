package page

import (
	"log"
	"runtime"
	"strconv"

	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/ark-go/giostart/internal/split"
)

var leftList *[]layout.FlexChild
var rightList *[]layout.FlexChild
var rightList2 *[]layout.Widget
var splitFlex *split.SplitFlexCol

func (p *TmainPage) initMainList() {
	log.Println("initMainList")
	splitFlex = &split.SplitFlexCol{}

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

								// dims := p.Split.Layout(gtx,
								// 	func(gtx C) D {
								// 		return layout.N.Layout(gtx, material.H3(p.Th, "Лево\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
								// 	}, func(gtx C) D {
								// 		return layout.N.Layout(gtx, material.H3(p.Th, "Право\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
								// 		//?  ///////////// 1вариант
								// 		// gtx.Constraints.Max.Y = 600
								// 		// //gtx.Constraints.Max = gtx.Constraints.Min
								// 		// dims2 := p.SplitV.Layout(gtx,
								// 		// 	func(gtx C) D {

								// 		// 		log.Println("ttts:", gtx.Constraints)

								// 		// 		return border.Layout(gtx, material.H3(p.Th, "Верх\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
								// 		// 		//return layout.N.Layout(gtx, material.H3(p.Th, "Верх\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
								// 		// 	}, func(gtx C) D {
								// 		// 		return layout.N.Layout(gtx, material.H3(p.Th, "Низ\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)

								// 		// 	})
								// 		// return dims2
								// 		//?  /////////// 2 Вариант
								// 		//return layout.Rigid(func(gtx C) D {}
								// 		//?  //////////
								// 	})
								splitFlex.Init(40, 5)
								dims := splitFlex.Layout(gtx, leftList, rightList2)

								return dims //D{Size: image.Pt(50, 100)}
							}),
						)
					},
				),
			)
		},
	}

	p.MainList = mainWidget

	leftList = &[]layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			//	gtx.Constraints.Min.X = 200
			return border.Layout(gtx, material.H3(p.Th, "left\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
		}),
	}
	rightList = &[]layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return border.Layout(gtx, material.H3(p.Th, "right\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Center.Layout(gtx, material.H3(p.Th, "Выравниваем\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Center.Layout(gtx, // это распологает виджет относительно всего доступного родительского виджета
				func(gtx C) D {
					return border.Layout(gtx, material.H3(p.Th, "Процессор\nВот:  "+strconv.Itoa(runtime.NumCPU())+" Почемуто не выравнивается а должен?").Layout)
				})
		}),
	}
	rightList2 = &[]layout.Widget{
		func(gtx C) D {
			return border.Layout(gtx, material.H3(p.Th, "right\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
		},

		func(gtx C) D {
			return border.Layout(gtx, material.H3(p.Th, "right\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
		},
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Center.Layout(gtx, material.H3(p.Th, "Выравниваем\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
		},
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Center.Layout(gtx, // это распологает виджет относительно всего доступного родительского виджета
				func(gtx C) D {
					return border.Layout(gtx, material.H3(p.Th, "Процессор\nВот:  "+strconv.Itoa(runtime.NumCPU())+" Почемуто не выравнивается а должен?").Layout)
				})
		},
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Center.Layout(gtx, // это распологает виджет относительно всего доступного родительского виджета
				func(gtx C) D {
					return border.Layout(gtx, material.H3(p.Th, "Процессор\nВот:  "+strconv.Itoa(runtime.NumCPU())+" Почемуто не выравнивается а должен?").Layout)
				})
		},
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Center.Layout(gtx, // это распологает виджет относительно всего доступного родительского виджета
				func(gtx C) D {
					return border.Layout(gtx, material.H3(p.Th, "Процессор\nВот:  "+strconv.Itoa(runtime.NumCPU())+" Почемуто не выравнивается а должен?").Layout)
				})
		},
	}
	// rightList2 = &[]layout.FlexChild{
	// 	layout.Rigid(func(gtx C) D {
	// 		return border.Layout(gtx, material.H3(p.Th, "right\n"+gtx.Constraints.Min.String()+"\n"+gtx.Constraints.Max.String()).Layout)
	// 	}),
	// }

}

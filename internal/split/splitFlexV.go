package split

import (
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

type C = layout.Context
type D = layout.Dimensions

type SplitFlexV struct {
	// Соотношение сохраняет текущий макет.
	// 0 - центр, -1 полностью слева, 1 полностью справа.
	// для установки в процентах левого поля использовать Split.LeftPanelWidth
	Ratio float32
	// Полоса - это ширина для изменения размера макета.
	Bar            unit.Value
	BarColor       color.NRGBA
	BarCenterColor color.NRGBA
	// выступ ручки
	BarKnobLedg      int
	BarKnobLedgColor color.NRGBA
	drag             bool
	dragID           pointer.ID
	dragY            float32
	initVal          bool // если LeftPanelWidth используется в цикле отрисовки, выключает пересчет данных после инициализации
}

func (s *SplitFlexV) Init(leftProcent int, barWidth int) {
	if barWidth < 0 {
		barWidth = 0
	}
	s = &SplitFlexV{
		Bar:              unit.Dp(float32(barWidth)),
		BarCenterColor:   color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xFF},
		BarColor:         color.NRGBA{R: 0xA5, G: 0xA5, B: 0xA5, A: 0xFF},
		BarKnobLedg:      barWidth + 6,
		BarKnobLedgColor: color.NRGBA{R: 0xA5, G: 0xA5, B: 0xA5, A: 0xFF},
	}
	s.LeftPanelWidth(leftProcent)
}

// ширина левой панели в процентах 0.1-100
// устанавливает Split.Ratio
// используется только при инициализации компонента
func (s *SplitFlexV) LeftPanelWidth(p int) {
	if !s.initVal {
		if p > 0 && p <= 100 {
			s.Ratio = float32(1.0 / (100 / float32(p)))
			if s.Ratio == 0.5 {
				s.Ratio = 0
			} else if s.Ratio < 0.5 {
				s.Ratio = 1 - s.Ratio
				s.Ratio *= -1
			}
			s.initVal = true
		}
	}
}

// размер
//var defaultBarWidth = unit.Dp(5)

func (s *SplitFlexV) Layout(gtx layout.Context, childrenLeft, childrenRight []layout.FlexChild) layout.Dimensions {
	//? в gtx общий размер для двух панелей,
	proportion := (s.Ratio + 1) / 2
	//? высчитываем Max.X верхнего поля путем пропорции от основного и с вычетом бара перетаскивателя
	leftsize := int(proportion*float32(gtx.Constraints.Max.Y) - float32(gtx.Px(s.Bar)))

	rightoffset := leftsize + gtx.Px(s.Bar)
	rightsize := gtx.Constraints.Max.Y - rightoffset
	_ = rightsize
	{ // handle input
		// Избегайте воздействия на дерево ввода событиями указателя.
		stack := op.Save(gtx.Ops)

		for _, ev := range gtx.Events(s) {
			e, ok := ev.(pointer.Event)
			if !ok {
				continue
			}

			switch e.Type {
			case pointer.Press:
				if s.drag {
					break
				}

				s.dragID = e.PointerID

				s.dragY = e.Position.Y

			case pointer.Drag:
				if s.dragID != e.PointerID {
					break
				}
				if e.Position.Y < float32(gtx.Constraints.Min.Y+10) || e.Position.Y > float32(gtx.Constraints.Max.Y-10) {
					break // если перетаскиваем за пределы компонента
				}
				deltaY := e.Position.Y - s.dragY
				s.dragY = e.Position.Y

				deltaRatio := deltaY * 2 / float32(gtx.Constraints.Max.Y)
				s.Ratio += deltaRatio

			case pointer.Release:
				fallthrough
			case pointer.Cancel:
				s.drag = false
			}
		}
		stack.Load()
	}

	mainflex := layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Max.X = 50
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx, childrenLeft...)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx, childrenRight...)
		}),
	)
	_ = mainflex
	// var w1, w2 layout.Dimensions // поже замерим кто больше
	// {
	// 	// левый компонент
	// 	stack := op.Save(gtx.Ops)

	// 	gtx1 := gtx // делаем копию
	// 	gtx1.Constraints = layout.Exact(image.Pt(gtx.Constraints.Max.X, leftsize))
	// 	gtx1.Constraints.Max.Y = gtx.Constraints.Max.Y // берем из переданного
	// 	w1 = left(gtx1)                                // w1 нужно чтоб определить где больше право/лево
	// 	stack.Load()
	// }

	// {
	// 	// правый компонент
	// 	stack := op.Save(gtx.Ops)
	// 	gtx1 := gtx // копия
	// 	//gtx1.Constraints = layout.Exact(image.Pt(rightsize, gtx.Constraints.Max.Y))
	// 	gtx1.Constraints = layout.Exact(image.Pt(gtx.Constraints.Max.X, rightsize))
	// 	gtx1.Constraints.Max.Y = gtx.Constraints.Max.Y // задавали для виджета
	// 	op.Offset(f32.Pt(0, float32(rightoffset))).Add(gtx1.Ops)
	// 	w2 = right(gtx1) // w2 нужно чтоб определить где больше право/лево

	// 	stack.Load()
	// }
	// // высчитаем больший размер двух половин ,,???
	// gtxN := gtx
	// _ = w1.Size.Y + w2.Size.Y
	// log.Println("test:", w1.Size.Y, w2.Size.Y)
	// gtxN.Constraints.Max.Y = w1.Size.Y + w2.Size.Y

	// {
	// 	// Рисуем разделитель
	// 	stack := op.Save(gtx.Ops)

	// 	barRect := image.Rect(0, leftsize, gtxN.Constraints.Max.X, rightoffset)
	// 	// создадим clip чтоб закрасить цветом
	// 	clip.Rect{Min: barRect.Min, Max: barRect.Max}.Add(gtx.Ops)
	// 	paint.ColorOp{Color: s.BarColor}.Add(gtx.Ops)
	// 	paint.PaintOp{}.Add(gtx.Ops)
	// 	// определим область для мышки
	// 	pointer.Rect(barRect).Add(gtx.Ops)
	// 	pointer.InputOp{Tag: s,
	// 		Types: pointer.Press | pointer.Drag | pointer.Release, // только для этих событий
	// 		Grab:  s.drag,
	// 	}.Add(gtx.Ops)
	// 	//! так установим курсор для перетаскивания
	// 	pointer.CursorNameOp{Name: pointer.CursorRowResize}.Add(gtx.Ops)
	// 	stack.Load()
	// }

	return mainflex
}

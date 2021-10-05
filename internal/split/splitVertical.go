package split

import (
	"image"
	"image/color"
	"log"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type SplitV struct {
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

func CreateV(leftProcent int, barWidth int) *SplitV {
	if barWidth < 0 {
		barWidth = 0
	}
	s := &SplitV{
		Bar:              unit.Dp(float32(barWidth)),
		BarCenterColor:   color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xFF},
		BarColor:         color.NRGBA{R: 0xA5, G: 0xA5, B: 0xA5, A: 0xFF},
		BarKnobLedg:      barWidth + 6,
		BarKnobLedgColor: color.NRGBA{R: 0xA5, G: 0xA5, B: 0xA5, A: 0xFF},
	}
	s.LeftPanelWidth(leftProcent)
	return s
}

// ширина левой панели в процентах 0.1-100
// устанавливает Split.Ratio
// используется только при инициализации компонента
func (s *SplitV) LeftPanelWidth(p int) {
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

func (s *SplitV) Layout(gtx layout.Context, left, right layout.Widget) layout.Dimensions {
	//? в gtx общий размер для двух панелей,
	proportion := (s.Ratio + 1) / 2
	//? высчитываем Max.X верхнего поля путем пропорции от основного и с вычетом бара перетаскивателя
	leftsize := int(proportion*float32(gtx.Constraints.Max.Y) - float32(gtx.Px(s.Bar)))

	rightoffset := leftsize + gtx.Px(s.Bar)
	rightsize := gtx.Constraints.Max.Y - rightoffset

	{ // handle input
		// Avoid affecting the input tree with pointer events.
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

		// // register for input
		// barRect := image.Rect(leftsize, 0, rightoffset, gtx.Constraints.Max.Y)
		// // создадим clip чтоб закрасить цветом
		// clip.Rect{Min: barRect.Min, Max: barRect.Max}.Add(gtx.Ops)
		// paint.ColorOp{Color: color.NRGBA{R: 0xA5, G: 0xA5, B: 0xA5, A: 0xFF}}.Add(gtx.Ops)
		// paint.PaintOp{}.Add(gtx.Ops)
		// // определим область для мышки
		// pointer.Rect(barRect).Add(gtx.Ops)
		// pointer.InputOp{Tag: s,
		// 	Types: pointer.Press | pointer.Drag | pointer.Release, // только для этих событий
		// 	Grab:  s.drag,
		// }.Add(gtx.Ops)

		// test рисуем сверху еще чтото
		// barRect = image.Rect(leftsize+2, 0, rightoffset-2, gtx.Constraints.Max.Y)
		// // создадим clip чтоб закрасить цветом
		// clip.Rect{Min: barRect.Min, Max: barRect.Max}.Add(gtx.Ops)
		// paint.ColorOp{Color: color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xFF}}.Add(gtx.Ops)
		// paint.PaintOp{}.Add(gtx.Ops)
		// end test
		// // test рисуем сверху еще чтото
		// p10 := 70 // gtx.Constraints.Max.Y * 10 / 100
		// t10 := (gtx.Constraints.Max.Y / 2) - (p10 / 2)
		// //barRect = image.Rect(leftsize-2, 0, rightoffset+2, gtx.Constraints.Max.Y)
		// log.Println(leftsize-5, t10, rightoffset+5, p10, gtx.Constraints.Max.Y)
		// barRect = image.Rect(leftsize-5, t10, rightoffset+5, p10)
		// // создадим clip чтоб закрасить цветом
		// clip.Rect{Min: barRect.Min, Max: barRect.Max}.Add(gtx.Ops)
		// paint.ColorOp{Color: color.NRGBA{R: 0xc5, G: 0xc5, B: 0xc5, A: 0xFF}}.Add(gtx.Ops)
		// paint.PaintOp{}.Add(gtx.Ops)
		// //! так установим курсор для перетаскивания
		//pointer.CursorNameOp{Name: pointer.CursorColResize}.Add(gtx.Ops)
		stack.Load()
	}
	var w1, w2 layout.Dimensions // поже замерим кто больше
	{
		// левый компонент
		stack := op.Save(gtx.Ops)

		gtx1 := gtx // делаем копию
		gtx1.Constraints = layout.Exact(image.Pt(gtx.Constraints.Max.X, leftsize))
		gtx1.Constraints.Max.Y = gtx.Constraints.Max.Y // берем из переданного
		w1 = left(gtx1)                                // w1 нужно чтоб определить где больше право/лево
		stack.Load()
	}

	{
		// правый компонент
		stack := op.Save(gtx.Ops)
		gtx1 := gtx // копия
		//gtx1.Constraints = layout.Exact(image.Pt(rightsize, gtx.Constraints.Max.Y))
		gtx1.Constraints = layout.Exact(image.Pt(gtx.Constraints.Max.X, rightsize))
		gtx1.Constraints.Max.Y = gtx.Constraints.Max.Y // задавали для виджета
		op.Offset(f32.Pt(0, float32(rightoffset))).Add(gtx1.Ops)
		w2 = right(gtx1) // w2 нужно чтоб определить где больше право/лево

		stack.Load()
	}
	// высчитаем больший размер двух половин ,,???
	gtxN := gtx
	_ = w1.Size.Y + w2.Size.Y
	log.Println("test:", w1.Size.Y, w2.Size.Y)
	gtxN.Constraints.Max.Y = w1.Size.Y + w2.Size.Y
	// if w1.Size.Y >= w2.Size.Y {
	// 	gtxN.Constraints.Max.Y = w1.Size.Y
	// } else {
	// 	gtxN.Constraints.Max.Y = w2.Size.Y
	// }

	{
		// Рисуем разделитель
		stack := op.Save(gtx.Ops)

		barRect := image.Rect(0, leftsize, gtxN.Constraints.Max.X, rightoffset)
		// создадим clip чтоб закрасить цветом
		clip.Rect{Min: barRect.Min, Max: barRect.Max}.Add(gtx.Ops)
		paint.ColorOp{Color: s.BarColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		// определим область для мышки
		pointer.Rect(barRect).Add(gtx.Ops)
		pointer.InputOp{Tag: s,
			Types: pointer.Press | pointer.Drag | pointer.Release, // только для этих событий
			Grab:  s.drag,
		}.Add(gtx.Ops)
		//! так установим курсор для перетаскивания
		pointer.CursorNameOp{Name: pointer.CursorRowResize}.Add(gtx.Ops)
		stack.Load()
	}

	{
		// if rightoffset-leftsize > 4 {
		// 	stack := op.Save(gtx.Ops)
		// 	// украшаем
		// 	barRect := image.Rect(leftsize+2, 0, rightoffset-2, gtxN.Constraints.Max.X)
		// 	// создадим clip чтоб закрасить цветом
		// 	clip.Rect{Min: barRect.Min, Max: barRect.Max}.Add(gtx.Ops)

		// 	paint.ColorOp{Color: s.BarCenterColor}.Add(gtx.Ops)
		// 	paint.PaintOp{}.Add(gtx.Ops)
		// 	stack.Load()
		// }
	}

	{
		// if s.BarKnobLedg > 0 {
		// 	//colorBar := color.NRGBA{R: 0xA5, G: 0xA5, B: 0xA5, A: 0xFF}
		// 	// рисуем ручку
		// 	stack := op.Save(gtx.Ops)
		// 	// test рисуем сверху еще чтото
		// 	p10 := gtxN.Constraints.Max.X * 10 / 100        // высота ручки 10 %
		// 	t10 := (gtxN.Constraints.Max.X / 2) - (p10 / 2) // начало отрисовки
		// 	//barRect = image.Rect(leftsize-2, 0, rightoffset+2, gtx.Constraints.Max.Y)
		// 	// ширина утолщения

		// 	barKnob := s.BarKnobLedg / 2                             // с одной стороны
		// 	widths := (rightoffset + barKnob) - (leftsize - barKnob) // общая ширина
		// 	radius := float32(widths / 2)

		// 	barRect := image.Rect(leftsize-barKnob, t10, rightoffset+barKnob, t10+p10)

		// 	he := barRect.Max.X - barRect.Min.X
		// 	if he < widths {
		// 		radius = float32(he / 2)
		// 	}
		// 	// создадим clip чтоб закрасить цветом
		// 	//clip.Rect{Min: barRect.Min, Max: barRect.Max}.Add(gtx.Ops)

		// 	r := f32.Rectangle{Max: layout.FPt(barRect.Max), Min: layout.FPt(barRect.Min)}
		// 	//	stack = op.Save(gtx.Ops)
		// 	clip.UniformRRect(r, radius).Add(gtx.Ops)
		// 	colorBarA := s.BarKnobLedgColor
		// 	colorBarA.A = 0xd0
		// 	paint.Fill(gtx.Ops, colorBarA)
		// 	//	stack.Load()
		// 	//////
		// 	// paint.FillShape(gtx.Ops,
		// 	// 	colorBar,
		// 	// 	clip.Stroke{
		// 	// 		Path:  cl.Path(gtx.Ops),
		// 	// 		Style: clip.StrokeStyle{Width: 1}, // рамку нарисует
		// 	// 		//	Style: clip.StrokeStyle{Width: width},
		// 	// 	}.Op(),
		// 	// )
		// 	//////

		// 	//paint.ColorOp{Color: color.NRGBA{R: 0x55, G: 0x55, B: 0x55, A: 0xFF}}.Add(gtx.Ops)
		// 	//paint.PaintOp{}.Add(gtx.Ops)
		// 	// stack = op.Save(gtx.Ops)
		// 	// stack.Load()
		// 	// Добавим область для перетаскивания
		// 	pointer.Rect(barRect).Add(gtx.Ops)
		// 	//pointer.AreaOp{}.Add(gtx.Ops) // --------------------------------------------------
		// 	pointer.InputOp{Tag: s,
		// 		Types: pointer.Press | pointer.Drag | pointer.Release, // только для этих событий
		// 		Grab:  s.drag,
		// 	}.Add(gtx.Ops)
		// 	//! так установим курсор для перетаскивания
		// 	pointer.CursorNameOp{Name: pointer.CursorColResize}.Add(gtx.Ops)
		// 	stack.Load()
		// }
	}
	//log.Println(gtx)
	//! вернем получившийся размер TODO: не забыть про прокрутку !!! т.е . ограничение надо расчитывать
	return layout.Dimensions{Size: gtxN.Constraints.Max}
}

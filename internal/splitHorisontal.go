package internal

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Split struct {
	// Соотношение сохраняет текущий макет.
	// 0 - центр, -1 полностью слева, 1 полностью справа.
	// для установки в процентах левого поля использовать Split.LeftPanelWidth
	Ratio float32
	// Полоса - это ширина для изменения размера макета.
	Bar unit.Value

	drag    bool
	dragID  pointer.ID
	dragX   float32
	initVal bool // если LeftPanelWidth используется в цикле отрисовки, выключает пересчет данных после инициализации
}

// ширина левой панели в процентах 0.1-100
// устанавливает Split.Ratio
// используется только при инициализации компонента
func (s *Split) LeftPanelWidth(p float32) {
	if !s.initVal {
		if p > 0 && p <= 100 {
			s.Ratio = float32(1.0 / (100 / p))
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
var defaultBarWidth = unit.Dp(5)

func (s *Split) Layout(gtx layout.Context, left, right layout.Widget) layout.Dimensions {

	//? в gtx общий размер для двух панелей, только Max.X  (Min.X такой же )
	bar := gtx.Px(s.Bar)
	if bar <= 1 {
		bar = gtx.Px(defaultBarWidth)
	}

	proportion := (s.Ratio + 1) / 2
	//? высчитываем Max.X левого поля путем пропорции от основного и с вычетом бара перетаскивателя
	leftsize := int(proportion*float32(gtx.Constraints.Max.X) - float32(bar))

	rightoffset := leftsize + bar
	rightsize := gtx.Constraints.Max.X - rightoffset

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

				s.dragX = e.Position.X

			case pointer.Drag:
				if s.dragID != e.PointerID {
					break
				}
				if e.Position.X < float32(gtx.Constraints.Min.X+10) || e.Position.X > float32(gtx.Constraints.Max.X-10) {
					break // если перетаскиваем за пределы компонента
				}
				deltaX := e.Position.X - s.dragX
				s.dragX = e.Position.X

				deltaRatio := deltaX * 2 / float32(gtx.Constraints.Max.X)
				s.Ratio += deltaRatio

			case pointer.Release:
				fallthrough
			case pointer.Cancel:
				s.drag = false
			}
		}

		// register for input
		barRect := image.Rect(leftsize, 0, rightoffset, gtx.Constraints.Max.Y)
		// создадим clip чтоб закрасить цветом
		clip.Rect{Min: barRect.Min, Max: barRect.Max}.Add(gtx.Ops)
		paint.ColorOp{Color: color.NRGBA{R: 0xA5, G: 0xA5, B: 0xA5, A: 0xFF}}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		// определим область для мышки
		pointer.Rect(barRect).Add(gtx.Ops)
		pointer.InputOp{Tag: s,
			Types: pointer.Press | pointer.Drag | pointer.Release, // только для этих событий
			Grab:  s.drag,
		}.Add(gtx.Ops)

		// test рисуем сверху еще чтото
		barRect = image.Rect(leftsize+2, 0, rightoffset-2, gtx.Constraints.Max.Y)
		// создадим clip чтоб закрасить цветом
		clip.Rect{Min: barRect.Min, Max: barRect.Max}.Add(gtx.Ops)
		paint.ColorOp{Color: color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xFF}}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		// end test

		//! так установим курсор для перетаскивания
		pointer.CursorNameOp{Name: pointer.CursorColResize}.Add(gtx.Ops)
		stack.Load()
	}
	var w1, w2 layout.Dimensions // поже замерим кто больше
	{
		// левый компонент
		stack := op.Save(gtx.Ops)

		gtx1 := gtx // делаем копию
		gtx1.Constraints = layout.Exact(image.Pt(leftsize, gtx.Constraints.Max.Y))
		gtx1.Constraints.Min.Y = gtx.Constraints.Min.Y // берем из переданного
		w1 = left(gtx1)                                // w1 нужно чтоб определить где больше право/лево
		stack.Load()
	}

	{
		// правый компонент
		stack := op.Save(gtx.Ops)
		gtx1 := gtx // копия
		gtx1.Constraints = layout.Exact(image.Pt(rightsize, gtx.Constraints.Max.Y))
		gtx1.Constraints.Min.Y = gtx.Constraints.Min.Y // задавали для виджета
		op.Offset(f32.Pt(float32(rightoffset), 0)).Add(gtx1.Ops)
		w2 = right(gtx1) // w2 нужно чтоб определить где больше право/лево

		stack.Load()
	}
	// высчитаем больший размер двух половин
	gtxN := gtx
	if w1.Size.Y >= w2.Size.Y {
		gtxN.Constraints.Max.Y = w1.Size.Y
	} else {
		gtxN.Constraints.Max.Y = w2.Size.Y
	}
	// вернем получившийся размер
	return layout.Dimensions{Size: gtxN.Constraints.Max}
}

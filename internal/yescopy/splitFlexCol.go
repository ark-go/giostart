//go:build ignore
// +build ignore

/**
* https://github.com/ark-go/giostart
*
*
 */
package split

import (
	"image"
	"image/color"
	"log"
	"time"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type C = layout.Context
type D = layout.Dimensions
type InsetPanel struct {
	Top, Right, Bottom, Left unit.Value
}
type SplitFlexCol struct {
	// структура отступов для всех сторон
	// если не задана берется из InsetDefault
	InsetPanel *InsetPanel
	// отступ для всех сторон по умолчанию
	InsetDefault unit.Value
	// Соотношение сохраняет текущий макет.
	// 0 - центр, -1 полностью слева, 1 полностью справа.
	// для установки в процентах левого поля использовать LeftPanelWidth
	Ratio float32
	//  ширина для изменения размера макета.
	BarSize        float32
	barSizeCurrent float32
	barSizeEnter   float32
	BarColor       color.NRGBA
	BarCenterColor color.NRGBA
	// выступ ручки
	BarKnobLedg      int
	BarKnobLedgColor color.NRGBA
	// фон левой панели
	BackgroundLeft color.NRGBA
	// фон правой панели
	BackgroundRight color.NRGBA
	drag            bool
	dragID          pointer.ID
	dragX           float32
	initVal         bool // чтоб инициализацию использовать один раз, если она в потоке
	//
	Color color.NRGBA
	// ?клавиша нажата
	pressed bool
	// размер левого поля
	leftSize int
	// процент заданный при инициализации
	startLeftSize int
	//-----------------
	startTime time.Time
	duration  time.Duration
	enter     bool
	leave     bool
}

func (s *SplitFlexCol) Init(leftProcent int, barWidth int) {
	if s.initVal {
		return
	}
	if barWidth < 0 {
		barWidth = 0
	}
	s.startTime = time.Now()
	s.duration = 300 * time.Millisecond
	s.BarSize = float32(barWidth) //gtx.Px(unit.Value(float32(barWidth)))
	s.barSizeCurrent = float32(barWidth)
	s.barSizeEnter = 6
	s.BarCenterColor = color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xFF}
	s.BarColor = color.NRGBA{R: 0xc3, G: 0xc3, B: 0xc3, A: 0xFF}
	s.BarKnobLedg = barWidth + 6
	s.BarKnobLedgColor = color.NRGBA{R: 0xA5, G: 0xA5, B: 0xA5, A: 0xFF}
	s.Color = color.NRGBA{R: 204, G: 204, B: 204, A: 255}
	s.BackgroundLeft = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	s.BackgroundRight = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	s.startLeftSize = leftProcent
	s.LeftPanelWidth(leftProcent)
	s.initVal = true
}

// ширина левой панели в процентах 1-100
// устанавливает Split.Ratio
// используется только при инициализации компонента
func (s *SplitFlexCol) LeftPanelWidth(p int) {
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

//	мне нужен был gtx :)
//		переводит int в пиксели текущего экрана, очень помогает!
func intToPx(gtx C, v int) int {
	return gtx.Px(unit.Dp(float32(v)))
}
func floatToPx(gtx C, v float32) int {
	return gtx.Px(unit.Dp(v))
}
func (s *SplitFlexCol) Layout(gtx layout.Context, childrenLeft, childrenRight *[]layout.FlexChild) layout.Dimensions {
	{
		stack := op.Save(gtx.Ops)
		for _, ev := range gtx.Events(s) {
			e, ok := ev.(pointer.Event)
			if !ok {
				continue
			}

			switch e.Type {
			case pointer.Press:
				// нажимаем кнопку на мыши
				if s.drag {
					break
				}
				s.pressed = true
				s.dragID = e.PointerID
				s.dragX = e.Position.X

				if e.Buttons == pointer.ButtonSecondary { // правая левая средняя
					log.Println("Правая")
					s.LeftPanelWidth(s.startLeftSize)
				}

			case pointer.Drag:
				// тащим мышь с нажатой кнопкой
				if s.dragID != e.PointerID {
					break
				}
				// if e.Position.X < float32(gtx.Constraints.Min.X+10) || e.Position.X > float32(gtx.Constraints.Max.X-10) {
				// 	break // если перетаскиваем за пределы компонента
				// }
				deltaX := e.Position.X - s.dragX
				s.dragX = e.Position.X

				deltaRatio := deltaX * 2 / float32(gtx.Constraints.Max.X)
				s.Ratio += deltaRatio
				log.Println("Ratio", s.Ratio)
			case pointer.Enter:
				//s.barSizeCurrent = s.barSizeEnter
				if !s.enter {
					s.startTime = time.Now()
					s.enter = true
					s.leave = false
				}
			case pointer.Leave:
				//s.barSizeCurrent = s.BarSize
				if !s.leave {
					s.startTime = time.Now()
					s.leave = true
					s.enter = false
				}
			case pointer.Release:
				// отпустили кнопку
				s.pressed = false
				fallthrough
			case pointer.Cancel:
				// что-то прервалось, может системой
				s.drag = false
			}
		}
		stack.Load()
	}

	{ // расчет ширины левого поля
		proportion := (s.Ratio + 1) / 2
		s.leftSize = int(proportion * float32(gtx.Constraints.Max.X))
		// проверка на границы слева - справа, здесь ограничения по краям
		if s.leftSize < 0 {
			s.leftSize = 0
		}
		// справа
		if s.leftSize+floatToPx(gtx, s.barSizeCurrent) > gtx.Constraints.Max.X {
			s.leftSize = gtx.Constraints.Max.X - floatToPx(gtx, s.barSizeCurrent)
		}
	}
	{ // Установка области отслеживания мышки, и нужных нам событий
		stack := op.Save(gtx.Ops)
		// Ограничьте область для событий указателя.
		pointer.Rect(image.Rect(s.leftSize, 0, s.leftSize+floatToPx(gtx, s.barSizeCurrent), gtx.Constraints.Min.Y)).Add(gtx.Ops)
		// На что реагировать
		pointer.InputOp{
			Tag:   s,                                                                              // просто индентификатор области события, не обязательно s
			Types: pointer.Press | pointer.Release | pointer.Drag | pointer.Enter | pointer.Leave, // какие события мышки хотим, лишних нам не надо
			Grab:  s.drag,                                                                         // что-то тянем
		}.Add(gtx.Ops)
		// так установим вид курсора для перетаскивания
		pointer.CursorNameOp{Name: pointer.CursorColResize}.Add(gtx.Ops)
		stack.Load()
	}
	//------------ этот блок будем отдавать ------------------------------------
	// Формируем наш виджет
	inset := layout.UniformInset(unit.Dp(10)) // равные отступы для inset
	if s.InsetPanel != nil {
		inset = layout.Inset(*s.InsetPanel) // если заданы свои отступы, установим их
	}
	gtxM := gtx // нам нужен общий размер всего виджета
	mainflex := layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx, // вся ширина
		// левая сторона
		layout.Rigid(func(gtx C) D {
			s.setBackground(gtx, s.BackgroundLeft) // закрашиваем фон чтоб не был прозрачным
			gtx.Constraints.Max.X = s.leftSize     // ограничиваем левый элемент
			gtx.Constraints.Min.X = s.leftSize     // этим растянем элемент вместе со свободным местом, например рамки.

			dims := inset.Layout(gtx, func(gtx C) D { // вставим отступы , получим размеры с отступом
				gtx.Constraints.Max.X = gtxM.Constraints.Max.X //* левому содержимому разрешаем рисоваться на всем поле и за разделитель
				// этим мы разрешаем или нет сжиматься содержимому в левом поле, или будем наезжать правым
				// ?здесь можно поиграть с цифрой - она определит максимально несжимаемую ширину
				return layout.Flex{ // вертикальный flex,
					Axis: layout.Vertical,
				}.Layout(gtx, *childrenLeft...) // вставляем наших деток
			})
			dims.Size.X = s.leftSize // тут мы скажем, на самом деле такой, даже если он и другой ширины
			// это позволит наезжать правому элементу,  правый элемент отсчитывает себя от этого поля
			return dims
		}),
		// разделитель, нас мало волнует сам виджет здесь,
		// нужно зарезервировать место, чтоб раздвинуть панели
		// а рисовать будем сами по координатам, TODO:
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Max.X = floatToPx(gtx, s.barSizeCurrent)
			return s.Layout2(gtx)
		}),
		// правая сторона
		layout.Rigid(func(gtx C) D {
			s.setBackground(gtx, s.BackgroundRight) // закрашиваем фон чтоб не был прозрачным
			return inset.Layout(gtx, func(gtx C) D {
				return layout.Flex{ // вертикальный flex
					Axis: layout.Vertical,
				}.Layout(gtx, *childrenRight...) // детки
			})
		}),
	)
	s.Layout4(gtx) // :)  да, можно рисовать и после
	return mainflex
}

// ----------------------  border widget ---------------------
// рисуем разделитель, как душе угодно
func (s *SplitFlexCol) Layout2(gtx C) D {

	dims := D{Size: gtx.Constraints.Max} // максимальная ширина разделителя
	// можно поменять цвет при нажатии - не рассматривал еще
	// col := s.Color
	// if s.pressed {
	// 	col = color.NRGBA{R: 0xFF, A: 0xFF}
	// }

	barRect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	// создадим clip чтоб закрасить цветом
	clip.Rect{Min: barRect.Min, Max: barRect.Max}.Add(gtx.Ops)
	paint.ColorOp{Color: s.BarColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	//---------------- //! анимация увеличение!
	if s.enter {
		elapsed := time.Since(s.startTime) // это заменяет конструкцию  time.Now().Sub(s.startTime)  длительность duration  ( time.Now() - s.startTime )
		progress := elapsed.Seconds() / s.duration.Seconds()
		k := s.barSizeEnter * float32(progress) // s.barSizeEnter (6) какую цифру мы хотим получить в конце анимации
		if k > 2.5 {                            // задержка на пролет мышки, не реагировать если мышка не задержалась дольше
			s.barSizeCurrent = s.BarSize + k
		}
		if progress < 1 {
			// Индикатор выполнения еще не закончил анимироваться.
			op.InvalidateOp{}.Add(gtx.Ops)
		} else {
			progress = 1
			s.enter = false
		}
	}
	//---------------- //! анимация уменьшение ушла мышка
	if s.leave {
		elapsed := time.Since(s.startTime) // это заменяет конструкцию  time.Now().Sub(s.startTime)
		progress := elapsed.Seconds() / s.duration.Seconds()
		k := s.barSizeEnter * float32(progress)
		//log.Println(":::::::::::::::::::::::::", k)
		if k > 2.5 { // задержка на пролет мышки без остановки
			if ((s.BarSize + s.barSizeEnter) - k) < s.barSizeCurrent { // реагируем только на уменьшение текущего размера
				s.barSizeCurrent = (s.BarSize + s.barSizeEnter) - k
			}
		}
		if progress < 1 {
			// Индикатор выполнения еще не закончил анимироваться.
			op.InvalidateOp{}.Add(gtx.Ops)
		} else {
			progress = 1
			s.leave = false
		}
	}
	//--------------
	return dims
}

func (s *SplitFlexCol) Layout4(gtx layout.Context) layout.Dimensions {

	dims := D{Size: gtx.Constraints.Max} // максимальная ширина разделителя
	// можно поменять цвет при нажатии - не рассматривал еще
	// col := s.Color
	// if s.pressed {
	// 	col = color.NRGBA{R: 0xFF, A: 0xFF}
	// }

	barRect := image.Rect(10, 0, 40, 50)
	// создадим clip чтоб закрасить цветом
	clip.Rect{Min: barRect.Min, Max: barRect.Max}.Add(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{R: 0xFF, A: 0xFF}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return dims
}

func (s *SplitFlexCol) setBackground(gtx C, col color.NRGBA) {
	barRect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Min.Y)
	clip.Rect{Min: barRect.Min, Max: barRect.Max}.Add(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	//return dims
}

//-------------------------- здесь рисуется рамка, она рисуется вокруг координат т.е. объект будет шире на толшину рамки
//defer op.Save(gtx.Ops).Load()
//dims := w
// sz := layout.FPt(dims.Size) // переводим в f32.Point
// r := f32.Rectangle{Max: sz}
// rr := float32(gtx.Px(s.CornerRadius))
// width := float32(gtx.Px(s.Width))
// sz.X -= width
// sz.Y -= width

// r = r.Add(f32.Point{X: width * 0.5, Y: width * 0.5})

// paint.FillShape(gtx.Ops,
// 	col,
// 	clip.Stroke{
// 		Path:  clip.UniformRRect(r, rr).Path(gtx.Ops),
// 		Style: clip.StrokeStyle{Width: width},
// 	}.Op(),
// )
// Рисуем разделитель
//stack := op.Save(gtx.Ops)

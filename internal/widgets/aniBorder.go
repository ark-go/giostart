package widgets

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

// рисуем рамку
type Border struct {
	Color        color.NRGBA
	CornerRadius unit.Value
	Width        unit.Value
	pressed      bool
}

func (b *Border) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	// Avoid affecting the input tree with pointer events.
	defer op.Save(gtx.Ops).Load()
	dims := w(gtx)
	//------------
	// here we loop through all the events associated with this button.
	for _, e := range gtx.Events(b) {
		if e, ok := e.(pointer.Event); ok {
			switch e.Type {
			case pointer.Press:
				b.pressed = true
			case pointer.Release:
				b.pressed = false
			}
		}
	}
	// Confine the area for pointer events.
	pointer.Rect(image.Rect(0, 0, dims.Size.X, dims.Size.Y)).Add(gtx.Ops)
	//pointer.Rect(image.Rectangle{dims.Size.X,dims.Size.Y}.Add(gtx.Ops)
	pointer.InputOp{
		Tag:   b,
		Types: pointer.Press | pointer.Release,
	}.Add(gtx.Ops)

	// Draw the button.
	col := b.Color
	if b.pressed {
		col = color.NRGBA{R: 0xFF, A: 0xFF}
	}
	return b.draw(gtx, w, col)
	//----------- end --------------------

}

func (b *Border) draw(gtx layout.Context, w layout.Widget, col color.NRGBA) layout.Dimensions {
	defer op.Save(gtx.Ops).Load()
	dims := w(gtx)
	sz := layout.FPt(dims.Size) // переводим в f32.Point
	r := f32.Rectangle{Max: sz}
	rr := float32(gtx.Px(b.CornerRadius))
	width := float32(gtx.Px(b.Width))
	sz.X -= width
	sz.Y -= width

	r = r.Add(f32.Point{X: width * 0.5, Y: width * 0.5})

	paint.FillShape(gtx.Ops,
		col,
		clip.Stroke{
			Path:  clip.UniformRRect(r, rr).Path(gtx.Ops),
			Style: clip.StrokeStyle{Width: width},
		}.Op(),
	)

	return dims
}

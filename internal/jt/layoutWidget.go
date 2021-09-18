package jt

import (
	"image"

	"gioui.org/layout"
)

func LayoutWidget(ctx layout.Context, width, height int) layout.Dimensions {
	return layout.Dimensions{
		Size: image.Point{
			X: width,
			Y: height,
		},
	}
}

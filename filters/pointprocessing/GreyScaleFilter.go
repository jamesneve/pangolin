package greyscalefilter

import (
	"image"
)

type GreyScaleFilter struct {
	Image image.Image
}

func NewGreyScaleFilter(image image.Image) GreyScaleFilter {
	greyScaleFilter := GreyScaleFilter{Image: image}
	return greyScaleFilter
}

func (f *GreyScaleFilter) GreyscaleFilter() image.Image {
	bounds := f.Image.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			gray.Set(x, y, f.Image.At(x, y))
		}
	}

	return gray
}

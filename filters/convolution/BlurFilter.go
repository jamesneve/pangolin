package blurfilter

import (
	"image"
	"image/color"
)

type BlurFilter struct {
	OriginalImage image.Image
	Width int
	Height int
}

func NewBlurFilter(image image.Image) BlurFilter {
	bounds := image.Bounds()
	return BlurFilter{image, bounds.Max.X, bounds.Max.Y}
}

// A box blur applied three times approximates a Gaussian blur with error < 0.1%
// Because the box blur can be computed as horizontal and vertical 1D sliding windows,
// it's much more efficient to do it this way
func (f *BlurFilter) GaussianBlur(radius int) image.Image {
	firstPassImage := f.applyBoxBlur(radius, f.OriginalImage)
	secondPassImage := f.applyBoxBlur(radius, firstPassImage)
	thirdPassImage := f.applyBoxBlur(radius, secondPassImage)
	return thirdPassImage
}

func (f *BlurFilter) BoxBlur(radius int) image.Image {
	return f.applyBoxBlur(radius, f.OriginalImage)
}

func (f *BlurFilter) applyBoxBlur(radius int, inputImage image.Image) image.Image {
	if f.Width < radius * 2 + 1 || f.Height < radius * 2 + 1 {
		return f.OriginalImage
	}

	horizontallyBlurredImage := f.horizontalBoxBlur(radius, inputImage)
	verticallyBlurredImage := f.verticalBoxBlur(radius, horizontallyBlurredImage)
	return verticallyBlurredImage
}

func (f *BlurFilter) horizontalBoxBlur(radius int, inputImage image.Image) image.Image {
	blurredImage := image.NewRGBA(image.Rect(0, 0, f.Width, f.Height))
	windowSize := uint32(radius * 2 + 1)

	sem := make(chan bool, f.Height)

	for y := 0; y < f.Height; y++ {
		go func(y int) {
			slidingWindow := make([]color.Color, windowSize)
			for x := range slidingWindow {
				slidingWindow[x] = inputImage.At(x, y)
			}

			r, g, b, a := f.initialWindowAverage(slidingWindow)
			blurredColor := color.RGBA{uint8(r / windowSize), uint8(g / windowSize), uint8(b / windowSize), uint8(a / windowSize)}

			for x := 0; x < f.Width; x++ {
				if x <= radius || x >= (f.Width - radius) {
					blurredImage.Set(x, y, blurredColor)
				} else {
					r, g, b, a = f.nextWindowAverage(r, g, b, a, inputImage.At(x + radius, y), inputImage.At(x - radius - 1, y))
					blurredColor = color.RGBA{uint8(r / windowSize), uint8(g / windowSize), uint8(b / windowSize), uint8(a / windowSize)}
					blurredImage.Set(x, y, blurredColor)
				}
			}
			sem <- true
		} (y)
	}

	for i := 0; i < f.Height; i++ { <-sem }

	return blurredImage
}

func (f *BlurFilter) verticalBoxBlur(radius int, inputImage image.Image) image.Image {
	blurredImage := image.NewRGBA(image.Rect(0, 0, f.Width, f.Height))
	windowSize := uint32(radius * 2 + 1)

	sem := make(chan bool, f.Height)

	for x := 0; x < f.Width; x++ {

		go func(x int) {
			slidingWindow := make([]color.Color, windowSize)
			for y := range slidingWindow {
				slidingWindow[y] = inputImage.At(x, y)
			}

			r, g, b, a := f.initialWindowAverage(slidingWindow)
			blurredColor := color.RGBA{uint8(r / windowSize), uint8(g / windowSize), uint8(b / windowSize), uint8(a / windowSize)}

			for y := 0; y < f.Height; y++ {
				if y <= radius || y >= (f.Height - radius) {
					blurredImage.Set(x, y, blurredColor)
				} else {
					r, g, b, a = f.nextWindowAverage(r, g, b, a, inputImage.At(x, y + radius), inputImage.At(x, y - radius - 1))
					blurredColor = color.RGBA{uint8(r / windowSize), uint8(g / windowSize), uint8(b / windowSize), uint8(a / windowSize)}
					blurredImage.Set(x, y, blurredColor)
				}
			}
			sem <- true
		} (x)
	}

	for i := 0; i < f.Width; i++ { <-sem }

	return blurredImage
}

func (f *BlurFilter) initialWindowAverage(window []color.Color) (uint32, uint32, uint32, uint32) {
	var ra, ga, ba, aa uint32
	for i := range window {
		r, g, b, a := window[i].RGBA()
		ra += r / 0x101
		ga += g / 0x101
		ba += b / 0x101
		aa += a / 0x101
	}

	return ra, ga, ba, aa
}

func (f *BlurFilter) nextWindowAverage(r, g, b, a uint32, nextColor, previousColor color.Color) (uint32, uint32, uint32, uint32) {
	rp, gp, bp, ap := previousColor.RGBA()
	rn, gn, bn, an := nextColor.RGBA()

	var ra, ga, ba, aa uint32

	ra = r - (rp / 0x101) + (rn / 0x101)
	ga = g - (gp / 0x101) + (gn / 0x101)
	ba = b - (bp / 0x101) + (bn / 0x101)
	aa = a - (ap / 0x101) + (an / 0x101)

	return ra, ga, ba, aa
}

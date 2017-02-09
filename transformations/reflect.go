package transformations

import (
	"image"
)

func HorizontalReflect(im image.Image) image.Image {
	bounds := im.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	reflectedImage := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x <= (width + 1)/2; x++ {
			reflectedImage.Set(width - x, y, im.At(x, y))
			reflectedImage.Set(x, y, im.At(width - x, y))
		}
	}

	return reflectedImage
}

func VerticalReflect(im image.Image) image.Image {
	bounds := im.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	reflectedImage := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y <= (height + 1)/2; y++ {
			reflectedImage.Set(x, height - y, im.At(x, y))
			reflectedImage.Set(x, y, im.At(x, height - y))
		}
	}

	return reflectedImage
}
package transformations

import "image"

func RotateAnticlockwise(im image.Image) image.Image {
	bounds := im.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	rotatedImage := image.NewRGBA(image.Rect(0, 0, height, width))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			rotatedImage.Set(y, width - x, im.At(x, y))
		}
	}

	return rotatedImage
}

func RotateClockwise(im image.Image) image.Image {
	bounds := im.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	rotatedImage := image.NewRGBA(image.Rect(0, 0, height, width))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			rotatedImage.Set(height - y, x, im.At(x, y))
		}
	}

	return rotatedImage
}
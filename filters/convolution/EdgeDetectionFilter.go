package convolution

import (
	"image"
	"image/color"
	"math"

	"github.com/jamesneve/pangolin/filters/pointprocessing"
)

type EdgeDetectionFilter struct {
	OriginalImage image.Image
	Width int
	Height int
}

func NewEdgeDetectionFilter(image image.Image) EdgeDetectionFilter {
	bounds := image.Bounds()
	return EdgeDetectionFilter{image, bounds.Max.X, bounds.Max.Y}
}

func (f *EdgeDetectionFilter) SobelEdgeDetection() image.Image {
	greyScaleFilter := greyscalefilter.NewGreyScaleFilter(f.OriginalImage)
	greyScaleImage := greyScaleFilter.GreyscaleFilter()

	edgeDetectedImage := f.applySobelEdgeDetection(greyScaleImage)

	return edgeDetectedImage
}

func (f *EdgeDetectionFilter) applySobelEdgeDetection(greyscaleImage *image.Gray) image.Image {
	horizontallyXFilteredImage := f.applyXFilterHorizontally(greyscaleImage)
	verticallyXFilteredImage := f.applyXFilterVertically(horizontallyXFilteredImage)
	horizontallyYFilteredImage := f.applyYFilterHorizontally(greyscaleImage)
	verticallyYFilteredImage := f.applyYFilterVertically(horizontallyYFilteredImage)
	newImage := f.combineXYEdgeDetection(verticallyXFilteredImage, verticallyYFilteredImage)
	return newImage
}

func (f *EdgeDetectionFilter) applyXFilterVertically(inputImage *image.Gray) *image.Gray {
	newImage := image.NewGray(image.Rect(0, 0, f.Width, f.Height))

	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			if y < 1 || y > (f.Height - 2) {
				newImage.Set(x, y, inputImage.GrayAt(x, y))
			} else {
				floatPointValue := (float64)(inputImage.GrayAt(x, y - 1).Y) - (float64)(inputImage.GrayAt(x, y + 1).Y)
				newPointValue := (uint8)(math.Abs(floatPointValue))
				newImage.Set(x, y, color.Gray{newPointValue})
			}
		}
	}

	return newImage
}

func (f *EdgeDetectionFilter) applyXFilterHorizontally(inputImage *image.Gray) *image.Gray {
	newImage := image.NewGray(image.Rect(0, 0, f.Width, f.Height))

	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			if x < 1 || x > (f.Width - 2) {
				newImage.Set(x, y, inputImage.GrayAt(x, y))
			} else {
				newPointValue := inputImage.GrayAt(x - 1, y).Y + 2 * inputImage.GrayAt(x, y).Y + inputImage.GrayAt(x + 1, y).Y
				newImage.Set(x, y, color.Gray{newPointValue})
			}
		}
	}

	return newImage
}

func (f *EdgeDetectionFilter) applyYFilterVertically(inputImage *image.Gray) *image.Gray {
	newImage := image.NewGray(image.Rect(0, 0, f.Width, f.Height))

	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			if y < 1 || y > (f.Height - 2) {
				newImage.Set(x, y, inputImage.GrayAt(x, y))
			} else {
				newPointValue := inputImage.GrayAt(x, y - 1).Y + 2 * inputImage.GrayAt(x, y).Y + inputImage.GrayAt(x, y + 1).Y
				newImage.Set(x, y, color.Gray{newPointValue})
			}
		}
	}

	return newImage
}

func (f *EdgeDetectionFilter) applyYFilterHorizontally(inputImage *image.Gray) *image.Gray {
	newImage := image.NewGray(image.Rect(0, 0, f.Width, f.Height))

	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			if x < 1 || x > (f.Width - 2) {
				newImage.Set(x, y, inputImage.GrayAt(x, y))
			} else {
				floatPointValue := (float64)(inputImage.GrayAt(x - 1, y).Y) - (float64)(inputImage.GrayAt(x + 1, y).Y)
				newPointValue := (uint8)(math.Abs(floatPointValue))
				newImage.Set(x, y, color.Gray{newPointValue})
			}
		}
	}

	return newImage
}

func (f *EdgeDetectionFilter) combineXYEdgeDetection(xImage, yImage *image.Gray) *image.Gray {
	newImage := image.NewGray(image.Rect(0, 0, f.Width, f.Height))

	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Width; y++ {
			floatPointValue := (float64)(xImage.GrayAt(x, y).Y) * (float64)(xImage.GrayAt(x, y).Y) + (float64)(yImage.GrayAt(x, y).Y) * (float64)(yImage.GrayAt(x, y).Y)
			newPointValue := (uint8)(math.Sqrt(floatPointValue))
			newImage.Set(x, y, color.Gray{newPointValue})
		}
	}

	return newImage
}


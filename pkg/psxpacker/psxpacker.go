package psxpacker

import (
	"image"
)

type InputImage struct {
	Name  string
	Image image.Image
}

type PackedImage struct {
	Name    string
	SubRect image.Rectangle
}

type PackResult struct {
	Image image.Image
	Atlas []PackedImage
}

func Pack(direction Direction, width int, height int, images []InputImage) (PackResult, error) {
	packContainer := newPackContainer(direction)
	for _, image := range images {
		packContainer.insert(image)
	}

	return packContainer.pack(width, height)
}

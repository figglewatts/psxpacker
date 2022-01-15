package psxpacker

import (
	"fmt"
	"image"
	"image/draw"
	"sort"
)

type packContainer struct {
	images    map[int][]InputImage
	direction Direction
}

func newPackContainer(direction Direction) packContainer {
	return packContainer{
		make(map[int][]InputImage),
		direction,
	}
}

func (p *packContainer) insert(image InputImage) {
	dimension := p.direction.dimensionFromImage(image.Image)

	if dimensionImages, ok := p.images[dimension]; ok {
		// if slice exists, insert at the correct index (sorted by the opposite dimension)
		oppositeDimension := p.direction.opposite().dimensionFromImage(image.Image)

		// find the index to insert at
		insertionIndex := sort.Search(len(dimensionImages), func(i int) bool {
			imageDimension := p.direction.opposite().dimensionFromImage(dimensionImages[i].Image)

			// sort first by opposite dimension
			if imageDimension >= oppositeDimension {
				return true
			}

			// then by name
			return dimensionImages[i].Name <= image.Name
		})

		// insert
		p.images[dimension] = insertDimensionImageAt(dimensionImages, insertionIndex, image)
	} else {
		// otherwise if it didn't exist already, make a new slice with image in it
		p.images[dimension] = []InputImage{image}
	}
}

func (p *packContainer) pack(resultWidth int, resultHeight int) (PackResult, error) {
	atlas := []PackedImage{}
	resultDimensions := image.Pt(resultWidth, resultHeight)
	resultImage := image.NewRGBA(image.Rectangle{image.Pt(0, 0), resultDimensions})
	cursor := image.Pt(0, 0)

	// get the image dimensions so we can sort them and iterate
	i := 0
	imageDimensions := make([]int, len(p.images))
	for imageDimension := range p.images {
		imageDimensions[i] = imageDimension
		i++
	}
	sort.Sort(sort.Reverse(sort.IntSlice(imageDimensions)))

	// iterate dimensions and pack images
	for _, imageDimension := range imageDimensions {
		images := p.images[imageDimension]
		for _, imageToPack := range images {
			imageOppositeDimension := p.direction.opposite().dimensionFromImage(imageToPack.Image)

			// if this rect will go past the edge in opposite direction then loop in the direction
			cursorOppositeDimension := p.direction.opposite().dimensionFromPoint(cursor)
			resultOppositeDimension := p.direction.opposite().dimensionFromPoint(resultDimensions)
			if cursorOppositeDimension+imageOppositeDimension > resultOppositeDimension {
				// increase cursor by imageDimension in dimension
				cursor = cursor.Add(p.direction.exclusivePoint(image.Pt(imageDimension, imageDimension)))

				// reset the opposite direction by subtracting from itself
				cursor = cursor.Sub(p.direction.opposite().exclusivePoint(cursor))
			}

			// if we go off the direction edge, then we're unable to pack
			cursorDimension := p.direction.dimensionFromPoint(cursor)
			resultDimension := p.direction.dimensionFromPoint(resultDimensions)
			if cursorDimension+imageDimension > resultDimension {
				return PackResult{}, fmt.Errorf("unable to pack %v image at cursor %v in %v result",
					p.direction.imageDimensionString(imageToPack.Image),
					cursor,
					p.direction.imageDimensionString(resultImage))
			}

			// we now have all the data to pack this image
			packedImage := PackedImage{
				imageToPack.Name,
				image.Rectangle{cursor, cursor.Add(imageToPack.Image.Bounds().Size())},
			}
			atlas = append(atlas, packedImage)

			// draw it onto the result image
			draw.Draw(resultImage, packedImage.SubRect, imageToPack.Image, image.Point{}, draw.Src)

			// move cursor to next spot
			cursor = cursor.Add(p.direction.opposite().exclusivePoint(
				image.Pt(imageOppositeDimension, imageOppositeDimension)))
		}

		// we've finished this line of dimensions, move onto next and reset opposite direction
		cursor = cursor.Add(p.direction.exclusivePoint(image.Pt(imageDimension, imageDimension)))
		cursor = cursor.Sub(p.direction.opposite().exclusivePoint(cursor))
	}

	return PackResult{resultImage, atlas}, nil
}

func insertDimensionImageAt(data []InputImage, i int, image InputImage) []InputImage {
	if i == len(data) {
		// easy case, insert at end
		return append(data, image)
	}

	// make space by shifting values at insert index up one
	data = append(data[:i+1], data[i:]...)

	// insert new
	data[i] = image

	return data
}

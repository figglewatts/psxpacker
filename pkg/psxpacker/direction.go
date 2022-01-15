package psxpacker

import (
	"fmt"
	"image"
	"strings"
)

// a Direction of Height will make images be packed in horizontal rows of height sorted by width
// a Direction of Width will make images be packed in vertical columns of width sorted by height
type Direction int

const (
	Width Direction = iota
	Height
)

func DirectionFromString(str string) (Direction, error) {
	switch strings.ToLower(str) {
	case "width":
		return Width, nil
	case "height":
		return Height, nil
	default:
		return -1, fmt.Errorf("unknown direction '%v', must be one of ['width', 'height']", str)
	}
}

func (d Direction) dimensionFromImage(image image.Image) int {
	return d.dimensionFromRect(image.Bounds())
}

func (d Direction) dimensionFromPoint(point image.Point) int {
	switch d {
	case Width:
		return point.X
	default:
		return point.Y
	}
}

func (d Direction) exclusivePoint(point image.Point) image.Point {
	switch d {
	case Width:
		return image.Pt(point.X, 0)
	default:
		return image.Pt(0, point.Y)
	}
}

func (d Direction) dimensionFromRect(rect image.Rectangle) int {
	switch d {
	case Width:
		return rect.Dx()
	default:
		return rect.Dy()
	}
}

func (d Direction) opposite() Direction {
	switch d {
	case Width:
		return Height
	default:
		return Width
	}
}

func (d Direction) imageDimensionString(image image.Image) string {
	dimension := d.dimensionFromImage(image)
	oppositeDimension := d.opposite().dimensionFromImage(image)
	switch d {
	case Width:
		return fmt.Sprintf("%vx%v", dimension, oppositeDimension)
	default:
		return fmt.Sprintf("%vx%v", oppositeDimension, dimension)
	}
}

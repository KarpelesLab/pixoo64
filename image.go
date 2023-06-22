package pixoo64

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
)

type Image struct {
	// Pix holds the image's pixels, in R, G, B order.
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
	Time   int // image's display time in ms
}

// NewImage creates a new pixoo64 format image. Size can only be 16, 32 or 64 pixels. Time
// is expressed in milliseconds.
func NewImage(size, time int) (*Image, error) {
	switch size {
	case 16, 32, 64:
	default:
		return nil, errors.New("invalid image size")
	}

	img := &Image{
		Pix:    make([]uint8, size*size*3),
		Stride: size * 3,
		Rect:   image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{size, size}},
		Time:   time,
	}
	return img, nil
}

// ConvertImage will take any image.Image and return a suitable image object
// unless the image size is not one of 16, 32 or 64
// This is useful when loading for example a png image without having to care
// for its pixel format.
func ConvertImage(img image.Image, time int) (*Image, error) {
	if i, ok := img.(*Image); ok {
		return i, nil
	}
	siz := img.Bounds().Size()
	if siz.X != siz.Y {
		return nil, errors.New("invalid image size")
	}
	newImg, err := NewImage(siz.X, time)
	if err != nil {
		return nil, err
	}
	// draw img to newImg
	draw.Draw(newImg, image.Rectangle{Max: siz}, img, image.Point{}, draw.Over)

	return newImg, nil
}

func (i *Image) ColorModel() color.Model {
	// use rgba model but we actually ignore alpha
	return color.RGBAModel
}

func (i *Image) Bounds() image.Rectangle {
	return i.Rect
}

func (i *Image) PixOffset(x, y int) int {
	return (y-i.Rect.Min.Y)*i.Stride + (x-i.Rect.Min.X)*3
}

func (i *Image) At(x, y int) color.Color {
	return i.RGBAAt(x, y)
}

func (i *Image) RGBAAt(x, y int) color.RGBA {
	if !(image.Point{x, y}.In(i.Rect)) {
		return color.RGBA{}
	}
	pixOfft := i.PixOffset(x, y)
	pix := i.Pix[pixOfft : pixOfft+3 : pixOfft+3] // see https://golang.org/issue/27857
	return color.RGBA{R: pix[0], G: pix[1], B: pix[2], A: 0xff}
}

func (i *Image) Set(x, y int, c color.Color) {
	c1 := color.RGBAModel.Convert(c).(color.RGBA)
	i.SetRGBA(x, y, c1)
}

func (i *Image) SetRGBA(x, y int, c color.RGBA) {
	if !(image.Point{x, y}.In(i.Rect)) {
		return
	}
	pixOfft := i.PixOffset(x, y)
	pix := i.Pix[pixOfft : pixOfft+3 : pixOfft+3] // see https://golang.org/issue/27857
	pix[0] = c.R
	pix[1] = c.G
	pix[2] = c.B
}

func (i *Image) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(i.Rect)
	if r.Empty() {
		return &Image{}
	}
	pixOfft := i.PixOffset(r.Min.X, r.Min.Y)

	return &Image{
		Pix:    i.Pix[pixOfft:],
		Stride: i.Stride,
		Rect:   r,
		Time:   i.Time,
	}
}

func (i *Image) SendTo(dev *Pixoo64) error {
	// send a single image means sending an animation with only this image
	a := Anim{}
	err := a.AppendFrame(i, -1)
	if err != nil {
		return err
	}
	return a.SendTo(dev)
}

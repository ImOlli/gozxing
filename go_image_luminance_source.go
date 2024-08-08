package gozxing

import (
	errors "golang.org/x/xerrors"
	"image"
)

func NewBinaryBitmapFromImage(img image.Image) (*BinaryBitmap, error) {
	src := NewLuminanceSourceFromImage(img)
	return NewBinaryBitmap(NewHybridBinarizer(src))
}

type GoImageLuminanceSource struct {
	*RGBLuminanceSource
}

const divisor = 4 * 0xffff

func NewLuminanceSourceFromImage(img image.Image) LuminanceSource {
	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y

	luminance := make([]byte, width*height)
	index := 0
	// Optimize special cases.
	switch img := img.(type) {
	case *image.Gray:
		pix := img.Pix
		for i := range pix {
			luminance[i] = pix[i]
		}
	case image.RGBA64Image:
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				p := img.(*image.RGBA)
				i := p.PixOffset(x, y)
				s := p.Pix[i : i+4 : i+4]

				r := uint32(s[0])
				g := uint32(s[1])
				b := uint32(s[2])
				a := uint32(s[3])
				r = (r << 8) | r
				g = (g << 8) | g
				b = (b << 8) | b
				a = (a << 8) | a

				lum := (r + 2*g + b) * 255 / divisor
				luminance[index] = byte((lum*a + (0xffff-a)*255) / 0xffff)
				index++
			}
		}
	default:
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				p := img.(*image.RGBA)
				i := p.PixOffset(x, y)
				s := p.Pix[i : i+4 : i+4]
				
				r := uint32(s[0])
				g := uint32(s[1])
				b := uint32(s[2])
				a := uint32(s[3])
				r = (r << 8) | r
				g = (g << 8) | g
				b = (b << 8) | b
				a = (a << 8) | a

				lum := (r + 2*g + b) * 255 / divisor
				luminance[index] = byte((lum*a + (0xffff-a)*255) / 0xffff)
				index++
			}
		}
	}

	return &GoImageLuminanceSource{&RGBLuminanceSource{
		LuminanceSourceBase{width, height},
		luminance,
		width,
		height,
		0,
		0,
	}}
}

func (this *GoImageLuminanceSource) Crop(left, top, width, height int) (LuminanceSource, error) {
	cropped, e := this.RGBLuminanceSource.Crop(left, top, width, height)
	if e != nil {
		return nil, e
	}
	return &GoImageLuminanceSource{cropped.(*RGBLuminanceSource)}, nil
}

func (this *GoImageLuminanceSource) Invert() LuminanceSource {
	return LuminanceSourceInvert(this)
}

func (this *GoImageLuminanceSource) IsRotateSupported() bool {
	return true
}

func (this *GoImageLuminanceSource) RotateCounterClockwise() (LuminanceSource, error) {
	width := this.GetWidth()
	height := this.GetHeight()
	top := this.top
	left := this.left
	dataWidth := this.dataWidth
	oldLuminas := this.RGBLuminanceSource.luminances
	newLuminas := make([]byte, width*height)

	for j := 0; j < width; j++ {
		x := left + width - 1 - j
		for i := 0; i < height; i++ {
			y := top + i
			newLuminas[j*height+i] = oldLuminas[y*dataWidth+x]
		}
	}
	return &GoImageLuminanceSource{&RGBLuminanceSource{
		LuminanceSourceBase{height, width},
		newLuminas,
		height,
		width,
		0,
		0,
	}}, nil
}

func (this *GoImageLuminanceSource) RotateCounterClockwise45() (LuminanceSource, error) {
	return nil, errors.New("RotateCounterClockwise45 is not implemented")
}
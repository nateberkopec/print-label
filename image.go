package printlabel

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
)

func LoadImage(r io.Reader, tapeWidth float64) (*image.Gray, error) {
	info, ok := TapeInfos[tapeWidth]
	if !ok {
		return nil, fmt.Errorf("unsupported tape width: %gmm", tapeWidth)
	}
	src, err := png.Decode(r)
	if err != nil {
		return nil, err
	}
	bounds := src.Bounds()
	if bounds.Dy() != info.PrintablePixels {
		return nil, fmt.Errorf("image height is %dpx; %gmm tape needs exactly %dpx", bounds.Dy(), tapeWidth, info.PrintablePixels)
	}
	img := image.NewGray(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := range bounds.Dy() {
		for x := range bounds.Dx() {
			img.SetGray(x, y, printedGray(src.At(bounds.Min.X+x, bounds.Min.Y+y)))
		}
	}
	return img, nil
}

func printedGray(c color.Color) color.Gray {
	_, _, _, a := c.RGBA()
	if a == 0xffff && color.GrayModel.Convert(c).(color.Gray).Y == 0 {
		return color.Gray{Y: 0}
	}
	return color.Gray{Y: 255}
}

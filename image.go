package printlabel

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
)

func LoadImage(r io.ReadSeeker, tapeWidth float64) (*image.Gray, error) {
	info, ok := TapeInfos[tapeWidth]
	if !ok {
		return nil, fmt.Errorf("unsupported tape width: %gmm", tapeWidth)
	}
	config, err := png.DecodeConfig(r)
	if err != nil {
		return nil, err
	}
	if config.Height != info.PrintablePixels {
		return nil, fmt.Errorf("image height is %dpx; %gmm tape needs exactly %dpx", config.Height, tapeWidth, info.PrintablePixels)
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	src, err := png.Decode(r)
	if err != nil {
		return nil, err
	}
	bounds := src.Bounds()
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

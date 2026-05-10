package printlabel

import (
	"image"
	"image/png"
	"os"
)

func SavePreview(img image.Image, path string, scale int) error {
	if scale < 1 {
		scale = 1
	}
	bounds := img.Bounds()
	preview := image.NewGray(image.Rect(0, 0, bounds.Dx()*scale, bounds.Dy()*scale))
	for y := range preview.Bounds().Dy() {
		for x := range preview.Bounds().Dx() {
			preview.Set(x, y, img.At(bounds.Min.X+x/scale, bounds.Min.Y+y/scale))
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, preview)
}

func threshold(src *image.Gray) *image.Gray {
	out := image.NewGray(src.Bounds())
	for i, v := range src.Pix {
		if v < 128 {
			out.Pix[i] = 0
		} else {
			out.Pix[i] = 255
		}
	}
	return out
}

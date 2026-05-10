package printlabel

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

func RenderText(text string, cfg Config) (*image.Gray, error) {
	info, ok := TapeInfos[cfg.TapeWidth]
	if !ok {
		return nil, fmt.Errorf("unsupported tape width: %gmm", cfg.TapeWidth)
	}
	sfntFont, err := loadFont(cfg)
	if err != nil {
		return nil, err
	}
	fontSize, err := chooseFontSize(text, cfg, info, sfntFont)
	if err != nil {
		return nil, err
	}
	face, err := opentype.NewFace(sfntFont, faceOptions(fontSize))
	if err != nil {
		return nil, err
	}
	defer face.Close()
	return drawText(text, cfg, info, face), nil
}

func loadFont(cfg Config) (*sfnt.Font, error) {
	fontData, err := os.ReadFile(ExpandHome(cfg.Font))
	if err != nil {
		return nil, err
	}
	collection, err := sfnt.ParseCollection(fontData)
	if err != nil {
		return nil, err
	}
	return collection.Font(cfg.FontIndex)
}

func chooseFontSize(text string, cfg Config, info TapeInfo, sfntFont *sfnt.Font) (int, error) {
	if cfg.FontSize != "" && cfg.FontSize != "auto" {
		return strconv.Atoi(cfg.FontSize)
	}
	for size := info.PrintablePixels - 4; size > 6; size-- {
		face, err := opentype.NewFace(sfntFont, faceOptions(size))
		if err != nil {
			return 0, err
		}
		bounds, _ := (&font.Drawer{Face: face}).BoundString(text)
		face.Close()
		if fixedToCeil(bounds.Max.Y-bounds.Min.Y) <= info.PrintablePixels-4 {
			return size, nil
		}
	}
	return 7, nil
}

func drawText(text string, cfg Config, info TapeInfo, face font.Face) *image.Gray {
	d := &font.Drawer{Face: face}
	bounds, advance := d.BoundString(text)
	textWidth := fixedToCeil(advance)
	textHeight := fixedToCeil(bounds.Max.Y - bounds.Min.Y)
	img := image.NewGray(image.Rect(0, 0, textWidth+2*cfg.Margin, info.PrintablePixels))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)
	d = &font.Drawer{Dst: img, Src: image.Black, Face: face}
	d.Dot = fixed.Point26_6{X: fixed.I(cfg.Margin) - bounds.Min.X, Y: fixed.I((img.Bounds().Dy()-textHeight)/2) - bounds.Min.Y}
	d.DrawString(text)
	return threshold(img)
}

func faceOptions(size int) *opentype.FaceOptions {
	return &opentype.FaceOptions{Size: float64(size), DPI: 72, Hinting: font.HintingFull}
}

func fixedToCeil(v fixed.Int26_6) int {
	return int((v + 63) >> 6)
}

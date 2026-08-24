package printlabel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
)

func ImageToRaster(img *image.Gray, tapeWidth float64) ([][]byte, error) {
	info, ok := TapeInfos[tapeWidth]
	if !ok {
		return nil, fmt.Errorf("unsupported tape width: %gmm", tapeWidth)
	}
	raster := make([][]byte, img.Bounds().Dx())
	for x := 0; x < img.Bounds().Dx(); x++ {
		raster[x] = rasterLine(img, x, info)
	}
	return raster, nil
}

func rasterLine(img *image.Gray, x int, info TapeInfo) []byte {
	line := make([]byte, BytesPerLine)
	for y := 0; y < img.Bounds().Dy(); y++ {
		if color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y == 0 {
			pin := info.RightMarginPins + y
			line[pin/8] |= 1 << (7 - (pin % 8))
		}
	}
	return line
}

func writeRaster(buf *bytes.Buffer, rasterLines [][]byte) {
	for _, line := range rasterLines {
		buf.WriteByte(0x47) // Uppercase G per PT-P900W Raster Command Reference v1.02.
		_ = binary.Write(buf, binary.LittleEndian, uint16(BytesPerLine))
		buf.Write(line)
	}
}

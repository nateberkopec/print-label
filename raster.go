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
	blank := bytes.Repeat([]byte{0}, BytesPerLine)
	for _, line := range rasterLines {
		if bytes.Equal(line, blank) {
			buf.WriteByte(0x5a)
			continue
		}
		buf.WriteByte(0x47)
		_ = binary.Write(buf, binary.LittleEndian, uint16(BytesPerLine))
		buf.Write(line)
	}
}

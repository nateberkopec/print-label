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

func BuildPrintJob(rasterLines [][]byte, tapeWidth float64, marginDots int) ([]byte, error) {
	info, ok := TapeInfos[tapeWidth]
	if !ok {
		return nil, fmt.Errorf("unsupported tape width: %gmm", tapeWidth)
	}
	buf := bytes.NewBuffer(nil)
	writeHeader(buf, info, len(rasterLines), marginDots)
	writeRaster(buf, rasterLines)
	buf.WriteByte(0x1a)
	return buf.Bytes(), nil
}

func writeHeader(buf *bytes.Buffer, info TapeInfo, rasterCount, marginDots int) {
	buf.Write(bytes.Repeat([]byte{0}, 200))
	buf.Write([]byte{0x1b, 0x40, 0x1b, 0x69, 0x61, 0x01, 0x1b, 0x69, 0x4d, 0x40, 0x1b, 0x69, 0x4b, 0x08, 0x1b, 0x69, 0x64})
	_ = binary.Write(buf, binary.LittleEndian, uint16(marginDots))
	buf.Write([]byte{0x1b, 0x69, 0x7a, 0xc4, 0x00, info.MediaByte, 0x00})
	_ = binary.Write(buf, binary.LittleEndian, uint32(rasterCount))
	buf.Write([]byte{0x00, 0x00, 0x4d, 0x00})
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

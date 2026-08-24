package printlabel

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func BuildPrintJob(rasterLines [][]byte, tapeWidth float64, marginDots int) ([]byte, error) {
	info, ok := TapeInfos[tapeWidth]
	if !ok {
		return nil, fmt.Errorf("unsupported tape width: %gmm", tapeWidth)
	}
	buf := bytes.NewBuffer(nil)
	writeInitialization(buf)
	writePageHeader(buf, info, len(rasterLines), marginDots, 2)
	writeRaster(buf, rasterLines)
	buf.WriteByte(0x1a)
	return buf.Bytes(), nil
}

func BuildMultiPrintJob(pages [][][]byte, tapeWidth float64, marginDots int) ([]byte, error) {
	info, ok := TapeInfos[tapeWidth]
	if !ok {
		return nil, fmt.Errorf("unsupported tape width: %gmm", tapeWidth)
	}
	if len(pages) == 0 {
		return nil, fmt.Errorf("at least one label is required")
	}

	buf := bytes.NewBuffer(nil)
	writeInitialization(buf)
	for i, rasterLines := range pages {
		writePageHeader(buf, info, len(rasterLines), marginDots, pageRole(i, len(pages)))
		writeRaster(buf, rasterLines)
		if i == len(pages)-1 {
			buf.WriteByte(0x1a)
		} else {
			buf.WriteByte(0x0c)
		}
	}
	return buf.Bytes(), nil
}

func pageRole(index, count int) byte {
	if index == count-1 {
		return 2
	}
	if index == 0 {
		return 0
	}
	return 1
}

func writeInitialization(buf *bytes.Buffer) {
	buf.Write(bytes.Repeat([]byte{0}, 200))
	buf.Write([]byte{0x1b, 0x40})
}

func writePageHeader(buf *bytes.Buffer, info TapeInfo, rasterCount, marginDots int, pageRole byte) {
	buf.Write([]byte{0x1b, 0x69, 0x61, 0x01})
	writePrintInformation(buf, info, rasterCount, pageRole)
	buf.Write([]byte{0x1b, 0x69, 0x4d, 0x40})
	buf.Write([]byte{0x1b, 0x69, 0x41, 0x01})
	buf.Write([]byte{0x1b, 0x69, 0x4b, 0x08})
	buf.Write([]byte{0x1b, 0x69, 0x64})
	_ = binary.Write(buf, binary.LittleEndian, uint16(marginDots))
	buf.Write([]byte{0x4d, 0x00})
}

func writePrintInformation(buf *bytes.Buffer, info TapeInfo, rasterCount int, pageRole byte) {
	buf.Write([]byte{0x1b, 0x69, 0x7a, 0xc4, 0x00, info.MediaByte, 0x00})
	_ = binary.Write(buf, binary.LittleEndian, uint32(rasterCount))
	buf.Write([]byte{pageRole, 0x00})
}

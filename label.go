package printlabel

const (
	ConfigPath   = "~/.config/label/config.yaml"
	TotalPins    = 560
	BytesPerLine = 70
)

type TapeInfo struct {
	PrintablePixels int
	RightMarginPins int
	MediaByte       byte
}

var TapeInfos = map[float64]TapeInfo{
	3.5: {48, 264, 0x04},
	6:   {64, 256, 0x06},
	9:   {106, 235, 0x09},
	12:  {150, 213, 0x0c},
	18:  {234, 171, 0x12},
	24:  {320, 128, 0x18},
	36:  {454, 61, 0x24},
}

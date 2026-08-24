package printlabel

import (
	"fmt"
	"strings"
)

func (s Status) Err() error {
	var messages []string
	messages = appendFlags(messages, s.ErrorInformation1, []flagMessage{
		{0x01, "no media"},
		{0x02, "end of media"},
		{0x04, "cutter jam"},
		{0x08, "weak batteries"},
		{0x40, "high-voltage adapter error"},
	})
	messages = appendFlags(messages, s.ErrorInformation2, []flagMessage{
		{0x01, "wrong media"},
		{0x02, "expansion buffer full"},
		{0x04, "communication error"},
		{0x08, "communication buffer full"},
		{0x10, "cover open"},
		{0x20, "overheating"},
		{0x40, "black marking not detected"},
		{0x80, "system error"},
	})
	if message := extendedErrorMessage(s.ExtendedError); message != "" {
		messages = append(messages, message)
	}
	if len(messages) == 0 && s.StatusType == statusErrorOccurred {
		messages = append(messages, "unspecified printer error")
	}
	if len(messages) == 0 {
		return nil
	}
	return fmt.Errorf("printer error: %s", strings.Join(messages, ", "))
}

type flagMessage struct {
	mask    byte
	message string
}

func appendFlags(messages []string, value byte, flags []flagMessage) []string {
	for _, flag := range flags {
		if value&flag.mask != 0 {
			messages = append(messages, flag.message)
		}
	}
	return messages
}

func extendedErrorMessage(value byte) string {
	return map[byte]string{
		0x10: "FLe tape end",
		0x1d: "high-resolution or draft printing error",
		0x1e: "adapter connection error",
		0x21: "incompatible media",
	}[value]
}

package printlabel

func stripComment(line string) string {
	inQuote := rune(0)
	for i, r := range line {
		if (r == '\'' || r == '"') && inQuote == 0 {
			inQuote = r
			continue
		}
		if r == inQuote {
			inQuote = 0
			continue
		}
		if r == '#' && inQuote == 0 {
			return line[:i]
		}
	}
	return line
}

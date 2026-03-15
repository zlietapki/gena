package generator

import "strings"

func mergeLines(linesA []string, linesB []string) []string {
	linesA = trimEmptyLines(linesA)
	linesB = trimEmptyLines(linesB)

	header := blockHeader(linesA)
	if header == "" {
		header = blockHeader(linesB)
	}

	if header != "" {
		inner := dedupLines(stripBlock(linesA), stripBlock(linesB))
		return append([]string{header}, append(inner, ")")...)
	}

	return dedupLines(linesA, linesB)
}

func trimEmptyLines(lines []string) []string {
	out := make([]string, 0, len(lines))

	for _, line := range lines {
		//line = strings.TrimSpace(line)
		if strings.TrimSpace(line) != "" {
			out = append(out, line)
		}
	}

	return out
}

func blockHeader(lines []string) string {
	if len(lines) > 1 {
		if lines[0] == "require (" || lines[0] == "import (" {
			if lines[len(lines)-1] == ")" {
				return lines[0]
			}
		}
	}

	return ""
}

func stripBlock(lines []string) []string {
	if blockHeader(lines) != "" {
		return lines[1 : len(lines)-1]
	}

	return lines
}

func dedupLines(linesA []string, linesB []string) []string {
	seen := map[string]bool{}

	var result []string
	for _, line := range append(linesA, linesB...) {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}

	return result
}

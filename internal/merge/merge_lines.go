package merge

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
		return append([]string{header}, append(inner, ")\n")...)
	}

	return dedupLines(linesA, linesB)
}

func trimEmptyLines(lines []string) []string {
	i := 0
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}

	j := len(lines) - 1
	for j >= i && strings.TrimSpace(lines[j]) == "" {
		j--
	}

	return lines[i : j+1]
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

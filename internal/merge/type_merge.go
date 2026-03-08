package merge

import "strings"

func lines(result []string, b []string) []string {
	result = trimEmpty(result)
	b = trimEmpty(b)

	header := blockHeader(result)
	if header == "" {
		header = blockHeader(b)
	}

	if header != "" {
		inner := dedupLines(stripBlock(result), stripBlock(b))
		return append([]string{header}, append(inner, ")\n")...)
	}

	return dedupLines(result, b)
}

func trimEmpty(s []string) []string {
	i := 0
	for i < len(s) && strings.TrimSpace(s[i]) == "" {
		i++
	}
	j := len(s) - 1
	for j >= i && strings.TrimSpace(s[j]) == "" {
		j--
	}
	return s[i : j+1]
}

func blockHeader(s []string) string {
	if len(s) > 1 {
		if s[0] == "require (" || s[0] == "import (" {
			if s[len(s)-1] == ")" {
				return s[0]
			}
		}
	}
	return ""
}

func stripBlock(s []string) []string {
	if blockHeader(s) != "" {
		return s[1 : len(s)-1]
	}
	return s
}

func dedupLines(result, b []string) []string {
	seen := map[string]bool{}
	var data []string
	for _, line := range append(result, b...) {
		if !seen[line] {
			seen[line] = true
			data = append(data, line)
		}
	}
	return data
}

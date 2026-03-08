package merge

func lines(result []string, b []string) []string {
	header := blockHeader(result)
	if header == "" {
		header = blockHeader(b)
	}

	if header != "" {
		inner := dedupLines(stripBlock(result), stripBlock(b))
		return append([]string{header}, append(inner, ")")...)
	}

	return dedupLines(result, b)
}

func blockHeader(s []string) string {
	if len(s) >= 2 && s[len(s)-1] == ")" {
		if s[0] == "require (" || s[0] == "import (" {
			return s[0]
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

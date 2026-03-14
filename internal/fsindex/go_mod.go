package fsindex

import (
	"bufio"
	"os"
	"strings"

	"github.com/zlietapki/microboiler/internal/vfs"
)

func getGoModBlocks(path string) ([]vfs.Block, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var blocks []vfs.Block
	var inRequire bool
	var requireBuf []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "module ") {
			blocks = append(blocks, vfs.Block{Name: "module", Type: vfs.BlockTypeSingle, Data: []string{line}})
			continue
		}

		if isRegexp(line, `^go \d+\.\d+`) {
			blocks = append(blocks, vfs.Block{Name: "go_ver", Type: vfs.BlockTypeSingle, Data: []string{line}})
			continue
		}

		if line == "require (" {
			inRequire = true
			requireBuf = []string{line}
			continue
		}

		if inRequire {
			if line == ")" {
				inRequire = false
				var direct []string
				for _, s := range requireBuf[1:] {
					if !strings.HasSuffix(strings.TrimSpace(s), "// indirect") {
						direct = append(direct, s)
					}
				}
				if len(direct) > 0 {
					data := append([]string{requireBuf[0]}, direct...)
					data = append(data, line)
					blocks = append(blocks, vfs.Block{
						Name: "require",
						Type: vfs.BlockTypeMerge,
						Data: data,
					})
				}
				continue
			}
			requireBuf = append(requireBuf, line)
		}
	}

	return blocks, scanner.Err()
}

package indexchecker

import (
	"fmt"
)

func SameMode() error {
	fc, err := getCollector()
	if err != nil {
		return err
	}

	for path, entries := range fc.fileMap {
		if len(entries) < 2 {
			continue
		}

		checkSameFileMode(path, entries)
	}

	for path, entries := range fc.dirMap {
		if len(entries) < 2 {
			continue
		}

		checkSameDirMode(path, entries)
	}

	return nil
}

func checkSameFileMode(path string, entries []fileEntry) bool {
	ref := entries[0].file.Mode

	ok := true
	for _, e := range entries[1:] {
		if ref != e.file.Mode {
			fmt.Printf("File mode mismatch:\n\tpath=%q\n\t%s mode=%s\n\t%s mode=%s\n",
				path, entries[0].projName, ref, e.projName, e.file.Mode)
			ok = false
		}
	}

	return ok
}

func checkSameDirMode(path string, entries []dirEntry) bool {
	ref := entries[0].dir.Mode

	ok := true
	for _, e := range entries[1:] {
		if ref != e.dir.Mode {
			fmt.Printf("Dir mode mismatch:\n\tpath=%q\n\t%s mode=%s\n\t%s mode=%s\n",
				path, entries[0].projName, ref, e.projName, e.dir.Mode)
			ok = false
		}
	}

	return ok
}

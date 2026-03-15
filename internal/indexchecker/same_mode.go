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

		err = checkSameFileMode(path, entries)
		if err != nil {
			return err
		}
	}

	for path, entries := range fc.dirMap {
		if len(entries) < 2 {
			continue
		}

		err = checkSameDirMode(path, entries)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkSameFileMode(path string, entries []fileEntry) error {
	ref := entries[0].file.Mode

	for _, e := range entries[1:] {
		if ref != e.file.Mode {
			return fmt.Errorf("File mode mismatch:\n\tpath=%q\n\t%s mode=%s\n\t%s mode=%s\n",
				path, entries[0].projName, ref, e.projName, e.file.Mode)
		}
	}

	return nil
}

func checkSameDirMode(path string, entries []dirEntry) error {
	ref := entries[0].dir.Mode

	for _, e := range entries[1:] {
		if ref != e.dir.Mode {
			return fmt.Errorf("Dir mode mismatch:\n\tpath=%q\n\t%s mode=%s\n\t%s mode=%s\n",
				path, entries[0].projName, ref, e.projName, e.dir.Mode)
		}
	}

	return nil
}

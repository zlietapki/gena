package vfs

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func debug(msg string, args ...interface{}) {
	if Debug {
		fmt.Printf("DEBUG "+msg, args...)
	}
}

func ignoreFile(name string) bool {
	if name == "go.sum" {
		return true
	}

	return false
}

func ignoreDir(name string) bool {
	if strings.HasPrefix(name, ".") {
		return true
	}

	return false
}

func isTextFile(filePath string) bool {
	mtype, err := mimetype.DetectFile(filePath)
	if err != nil {
		return false
	}

	knownMimes := map[string]bool{
		"text/plain; charset=utf-8": true,
		"application/x-executable":  false,
	}

	if val, ok := knownMimes[mtype.String()]; ok {
		return val
	}

	fmt.Println("Unknown MIME type:", mtype.String(), " for file:", filePath)
	os.Exit(1)
	return false
}

func isRegexp(line string, reg string) bool {
	matched, err := regexp.MatchString(reg, line)
	if err != nil {
		panic(err)
	}

	return matched
}

package generate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

var Debug = false

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
	mime, err := mimetype.DetectFile(filePath)
	if err != nil {
		return false
	}

	knownMimes := map[string]bool{
		"text/plain; charset=utf-8": true,
		"application/x-executable":  false,
	}

	if val, ok := knownMimes[mime.String()]; ok {
		return val
	}

	fmt.Println("Warning: Unknown MIME type:", mime.String(), " for file:", filePath)

	return false
}

func isRegexp(line string, reg string) bool {
	matched, err := regexp.MatchString(reg, line)
	if err != nil {
		panic(err)
	}

	return matched
}

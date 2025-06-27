package core

import (
	"bytes"
	"io"
	"os"
	"unicode/utf8"
)

const maxReadBytes = 1024

// IsTextFile checks if the file appears to be textual (UTF-8, no nulls, few control chars)
func IsTextFile(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	buf := make([]byte, maxReadBytes)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}
	buf = buf[:n]

	// Null byte check (likely binary)
	if bytes.Contains(buf, []byte{0}) {
		return false, nil
	}

	// UTF-8 check
	if !utf8.Valid(buf) {
		return false, nil
	}

	// Control character ratio check (excluding \t, \n, \r)
	controlCount := 0
	for _, b := range buf {
		if b < 32 && b != 9 && b != 10 && b != 13 {
			controlCount++
		}
	}
	if len(buf) > 0 && float64(controlCount)/float64(len(buf)) > 0.3 {
		return false, nil
	}

	return true, nil
}

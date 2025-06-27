package core

import (
	"path/filepath"
	"testing"
)

// TestIsTextFile verifies the text detection logic by testing several known files:
// - valid UTF-8 text
// - binary file
// - empty file
// - invalid UTF-8 content
func TestIsTextFile(t *testing.T) {

	cases := []struct {
		name     string
		filename string
		want     bool
	}{
		{"valid text file", filepath.Join(BasePath, "sample.txt"), true},
		{"binary file", filepath.Join(BasePath, "image.jpg"), false},
		{"empty file", filepath.Join(BasePath, "empty.txt"), true},
		{"invalid utf8", filepath.Join(BasePath, "invalid_utf8.bin"), false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			abs, _ := filepath.Abs(c.filename)
			got, err := IsTextFile(abs)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.want {
				t.Errorf("IsTextFile(%q) = %v; want %v", c.filename, got, c.want)
			}
		})
	}
}

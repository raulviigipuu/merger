package core

import (
	"bytes"
	"path/filepath"
	"testing"
)

// TestCopyFileTo verifies that the contents of a known text file
// are correctly streamed into an in-memory buffer.
func TestCopyFileTo(t *testing.T) {
	var buf bytes.Buffer
	err := copyFileTo(&buf, filepath.Join(BasePath, "sample.txt"))
	if err != nil {
		t.Fatalf("copyFileTo failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Errorf("Expected content, got empty buffer")
	}
}

// TestCollectTextFilesRecursive verifies that text files are collected
// correctly from a nested directory structure, and no unexpected files are included.
func TestCollectTextFilesRecursive(t *testing.T) {
	var files []string
	root := filepath.Join(BasePath, "nested")
	err := collectTextFilesRecursive(root, &files)
	if err != nil {
		t.Fatalf("collectTextFilesRecursive failed: %v", err)
	}

	expected := map[string]bool{
		filepath.Join(root, "sub.md"): true,
	}

	for _, f := range files {
		if !expected[f] {
			t.Errorf("Unexpected file collected: %s", f)
		}
		delete(expected, f)
	}

	for f := range expected {
		t.Errorf("Expected file missing: %s", f)
	}
}

package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Run is the entry point for merging text files from inputPaths into a single outputPath.
// It collects valid text files, adds a header with their relative path, and appends their content.
func Run(inputPaths []string, outputPath string) error {
	fmt.Println("Output file:", outputPath)

	var allFiles []string

	// Collect all valid text files from each input path (file or directory).
	for _, path := range inputPaths {
		files, err := collectTextFiles(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping %s: %v\n", path, err)
			continue
		}
		// Append discovered files to the master list.
		allFiles = append(allFiles, files...)
	}

	// Abort if no text files were found.
	if len(allFiles) == 0 {
		return fmt.Errorf("no textual files found in input paths")
	}

	// Open the output file for writing (overwrites if it exists).
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Process each file: write header and content to output.
	for _, file := range allFiles {
		// Convert file path to relative format for clean display in output.
		relPath, err := filepath.Rel(".", file)
		if err != nil {
			relPath = file // fallback to absolute if conversion fails
		}
		fmt.Println("Merging:", relPath)

		// Write a header to identify which file content follows.
		header := fmt.Sprintf("==== %s ====\n", relPath)
		if _, err := outFile.WriteString(header); err != nil {
			return err
		}

		// Append the file's content to the output file.
		err = copyFileTo(outFile, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read %s: %v\n", file, err)
			continue
		}

		// Add padding (2 newlines) between files.
		if _, err := outFile.WriteString("\n\n"); err != nil {
			return err
		}
	}

	fmt.Println("Done.")
	return nil
}

// collectTextFiles handles both files and directories.
// If the path is a single text file, it checks and returns it.
// If the path is a directory, it recurses and collects all valid text files.
func collectTextFiles(path string) ([]string, error) {
	var files []string

	// Get info on the path to distinguish between file vs directory.
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat error for path %q: %w", path, err)
	}

	// If it's a single file, check if it's textual.
	if !info.IsDir() {
		ok, err := IsTextFile(path)
		if err == nil && ok {
			files = append(files, path)
		}
		return files, nil
	}

	// If it's a directory, traverse recursively.
	err = collectTextFilesRecursive(path, &files)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// collectTextFilesRecursive is a helper that walks through directories recursively,
// appending all discovered text file paths to the output slice pointer.
func collectTextFilesRecursive(dir string, out *[]string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir error for %q: %w", dir, err)
	}

	for _, entry := range entries {
		childPath := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			// Recurse into subdirectory.
			if err := collectTextFilesRecursive(childPath, out); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: skipping subdir %q: %v\n", childPath, err)
			}
			continue
		}

		// Check if file is textual and add to list if so.
		ok, err := IsTextFile(childPath)
		if err == nil && ok {
			*out = append(*out, childPath)
		}
	}

	return nil
}

// copyFileTo reads the contents of the file at 'path' and writes it to 'dst'.
// Used to stream each text file into the final output.
func copyFileTo(dst io.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer file.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		return fmt.Errorf("failed to copy content from %q: %w", path, err)
	}

	if written == 0 {
		fmt.Fprintf(os.Stderr, "Warning: file %q is empty\n", path)
	}

	return nil
}

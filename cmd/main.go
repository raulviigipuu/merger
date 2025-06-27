package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	merger "github.com/raulviigipuu/merger/internal/core"
)

var (
	Version  = "dev" // value from -ldflags
	output   string
	showHelp bool
	showVer  bool
)

// Special function that runs before main
func init() {
	flag.StringVar(&output, "o", "out.txt", "Output file path (default: out.txt)")
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.BoolVar(&showVer, "v", false, "Show version")
}

func main() {

	// Look for misplaced flags: if any flag appears *after* a non-flag argument
	for i, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") {
			continue // valid positional
		}
		if i > 0 && !strings.HasPrefix(os.Args[i], "-") {
			fmt.Fprintln(os.Stderr, "Error: flags must come before input files or directories.")
			fmt.Fprintln(os.Stderr, "Correct usage: merger -o out.txt file1 dir2 ...")
			os.Exit(1)
		}
	}

	flag.Parse()

	if showHelp {
		fmt.Fprintln(os.Stdout, "Usage: merger [OPTIONS] FILES and DIRECTORIES...")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Merges the content of textual files into a single output file.")
		fmt.Fprintln(os.Stdout, "Supports recursive directory traversal, path deduplication, and text-only filtering.")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Example:")
		fmt.Fprintln(os.Stdout, "  merger -o out.txt notes.txt docs/ src/file.md")
		fmt.Fprintln(os.Stdout)
		os.Exit(0)
	}

	if showVer {
		fmt.Printf("merger version: %s\n", Version)
		os.Exit(0)
	}

	inputs := flag.Args()
	if len(inputs) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no input files or directories specified.")
		fmt.Fprintln(os.Stderr, "Use -h for help.")
		os.Exit(1)
	}

	// Deduplicate and normalize input paths
	unique := make(map[string]struct{})
	var finalInputs []string
	for _, p := range inputs {
		abs, err := filepath.Abs(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping invalid path %q: %v\n", p, err)
			continue
		}
		if _, seen := unique[abs]; !seen {
			unique[abs] = struct{}{}
			finalInputs = append(finalInputs, abs)
		}
	}

	if err := merger.Run(finalInputs, output); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

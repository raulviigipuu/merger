# merger

A minimal, fast CLI tool that merges the contents of multiple **text files** into a single `.txt` output file. Supports recursive directory traversal and outputs clean, annotated text.

---

## ✨ Features

- Accepts both files and directories as input
- Recursively searches directories for textual files
- Skips binary files using heuristics (null byte, UTF-8 validation, control character ratio)
- Annotates each file in the output with a header containing its relative path
- Deduplicates input paths
- Default output file is `out.txt` (can be overridden)
- Logs all actions to stdout/stderr

---

## 🔧 Usage

```
merger [OPTIONS] [FILES and DIRECTORIES...]
```

### Example

```
merger notes.txt ./src ./docs 
```

```
merger -o merged.txt notes.txt ./src ./docs 
```

This will:
- Search `notes.txt`, `./src/`, and `./docs/`
- Collect all valid **textual files**
- Write their contents to `merged.txt`
- Each file is preceded by a line like:

```
==== path/to/file.txt ====
<file contents>
```

---

## 🛠 Options

| Flag          | Description                           |
|---------------|---------------------------------------|
| `-o`          | Output file path (default: `out.txt`) |
| `-v`          | Print version and exit                |
| `-h`          | Show help message                     |

---

## 💡 Text Detection Heuristics

A file is considered *textual* if:

- It contains **no null bytes**
- It passes **UTF-8 validation**
- It has **<30% control characters** (excluding tab, newline, carriage return)

---

## 🚧 Future Plans

This tool is intentionally minimal — but the following features are planned:

- 📦 **Binary file support via base64 encoding**  
  Include non-textual files in the merged output in a reversible, printable form.

- 🔄 **Unpacking mode**  
  Reconstruct original directory structure and files from a merged text file.

---

## 📦 Building from Source

### Build locally:

```
go build -ldflags "-X main.Version=1.0.0" -o merger ./cmd/merger
```

### Install locally (to GOBIN)

```
go install -ldflags "-X main.Version=1.0.0" ./cmd/merger
```

### Run:

```
./merger -h
```

---

## 🔧 Maintenance

Upgrading go version in go.mod

```
go mod edit -go 1.25.3

go mod tidy
```

---

## 🧪 Running Tests

To run the unit tests:

```sh
go test ./internal/...

go test -v ./internal/...
```

## 📄 License

MIT

---
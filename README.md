# Concurrent Text Analyzer — Project Specification

## Goal

Build a Go command-line program that scans a folder of text files, counts words, and prints a ranked report. It should continue processing when one file fails and should support processing several files concurrently.

This project is complete when the same input produces the same word counts with one worker or several workers.

## How a user runs it

```bash
go run . -dir ./samples -workers 4 -top 10
```

| Option | Meaning | Default |
| --- | --- | --- |
| `-dir` | Folder to scan, including its subfolders | `.` |
| `-workers` | Maximum number of files processed at once | `4` |
| `-top` | Number of words to show in the ranking | `10` |

The program should reject an empty directory path and values of `-workers` or `-top` below `1` with a clear error message.

## Required functionality

### 1. Discover files

- Walk the requested folder recursively.
- Include regular files whose extension is `.txt` (case-insensitive, so `.TXT` also counts).
- Ignore other file types and directories.
- Count how many eligible files were found.
- If the starting folder does not exist or cannot be scanned, report the problem and exit with an error.
- If no text files are found, print a valid report with zero counts.

### 2. Define and count words

- Treat a word as a consecutive sequence of Unicode letters.
- Treat spaces, numbers, and punctuation as separators. For example, `go123lang` becomes `go` and `lang`.
- Convert letters to lowercase before counting. `Go`, `GO`, and `go` are the same word.
- Count each occurrence, not just each distinct word.
- Empty files contribute zero words and still count as successfully processed.

Example:

```text
Input:  Go, go! Café 123 café.
Counts: go=2, café=2
Total:  4 words
```

The word-counting function should accept an `io.Reader`. This lets it read a file in the program and a string in a test.

### 3. Process files concurrently

- Put discovered file paths into a jobs channel.
- Start the number of worker goroutines specified by `-workers`.
- Each worker reads one file at a time and produces a result containing the file path, its word counts, and any error.
- Collect all results before printing the final report.
- Use a `sync.WaitGroup` to know when all workers have finished.
- Protect shared totals and the combined word map with a `sync.Mutex`. Keep the mutex locked only while reading or changing that shared state.
- Do not let a failure in one file prevent other files from being processed.

### 4. Report results

Print:

1. Number of eligible files found.
2. Number processed successfully.
3. Number that failed.
4. Total number of words in successful files.
5. The requested number of most common words and their counts.
6. The path and error for each failed file, if any.

Sort words by count from highest to lowest. For words with the same count, sort alphabetically so the report is predictable.

Example format (the numbers are illustrative):

```text
Files found:     8
Files processed: 7
Files failed:    1
Total words:     1243

Top 10 words
1. go          82
2. program     47
3. file        39

Failed files
samples/broken.txt: permission denied
```

If there are fewer distinct words than `-top`, print only the words available. The order in which workers finish must not change the final ranking.

## Suggested project layout

```text
text-analyzer/
├── go.mod
├── main.go                 # flags, validation, final output
├── internal/
│   ├── analyzer/
│   │   ├── files.go        # find .txt files
│   │   ├── words.go        # tokenize and count one reader
│   │   ├── pool.go         # jobs, workers, and results
│   │   ├── stats.go        # shared totals and mutex
│   │   └── words_test.go   # word-counting tests
│   └── counter/
│       └── counter.go      # generic Counter[T comparable]
└── samples/
    ├── one.txt
    └── two.txt
```

Create the module with `go mod init example.com/text-analyzer`. The `example.com` name is suitable for a local exercise; replace it with your repository path if you publish the code.

## Where each Go topic fits

| Topic | Application |
| --- | --- |
| Variables and constants | Store flag values and define defaults or status values. |
| Formatting | Use `fmt` for errors and the final report. |
| Conditionals | Validate inputs and handle empty folders or failed files. |
| Functions | Separate discovery, counting, merging, sorting, and printing. |
| Structs | Model `Config`, `FileResult`, `Stats`, and a ranked `WordCount`. |
| Interfaces | Let the counting function accept `io.Reader`. |
| Errors | Return and wrap errors; continue after individual file failures. |
| Loops | Walk results, run workers, merge counts, and print rankings. |
| Slices | Store paths, results, and sorted word/count pairs. |
| Maps | Store counts by word. |
| Pointers | Pass a shared `*Stats` to workers. |
| Packages and modules | Organize the CLI and reusable code in one Go module. |
| Channels | Send jobs to workers and results back to the collector. |
| Mutexes | Protect the shared statistics while workers update them. |
| Generics | Make a reusable `Counter[T comparable]`; use `Counter[string]` for words. |
| Enum pattern | Define a `FileStatus` type with `iota` constants such as `Processed` and `Failed`. |

## Build order

1. **Single reader:** Implement and test word counting using short strings.
2. **One file:** Open a file, count it, and handle errors.
3. **Folder:** Find `.txt` files and process them one at a time.
4. **Report:** Merge counts, sort them, and print the requested top words.
5. **Concurrency:** Add the jobs channel, workers, results channel, wait group, and mutex.
6. **Generic counter:** Move reusable counting behavior into `Counter[T comparable]`.

Keep the program working after each step. Add concurrency only after the one-worker report is correct.

## Checks before calling it complete

- `Go, go! Café 123 café.` produces `go=2` and `café=2`.
- An empty text file succeeds and adds zero words.
- Files with other extensions are ignored.
- `-top 20` works even when there are fewer than 20 distinct words.
- Invalid flag values and a missing starting directory produce clear errors.
- A failed file is listed, while other files still contribute to the report.
- Running with `-workers 1` and `-workers 4` produces identical totals and rankings.
- `go fmt ./...`, `go test ./...`, and `go test -race ./...` complete successfully.

## Tools

- **Go toolchain:** compile, run, format, and test the program.
- **Code editor:** any editor with Go support.
- **Terminal:** run the program and its checks.
- **Go standard library:** `flag`, `fmt`, `io`, `os`, `path/filepath`, `sort`, `strings`, `sync`, `testing`, and `unicode`. No third-party package is required.
# FilesStats

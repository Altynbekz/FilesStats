# FilesStats: concurrent text analyzer

FilesStats is a Go exercise that scans a directory for `.txt` files, counts words in each file, and prints a separate report for each file. Several files can be processed at the same time with worker goroutines.

The code currently runs from `internal/analyzer`. The section [Project specification and remaining work](#project-specification-and-remaining-work) records the broader goal; it is not a list of features already finished.

## Run the current program

From the repository root, with a Go toolchain that supports the version in `go.mod`:

```bash
go run ./internal/analyzer -dir ./internal/test -workers 4 -top 10
```

You can also see the flags with `go run ./internal/analyzer -h`. Paths passed to `-dir` are relative to the **current working directory**, not to `internal/analyzer`. `go run .` from the repository root is not the analyzer entry point.

| Flag | Default | Current meaning |
| --- | --- | --- |
| `-dir` | `.` | Starting directory for the recursive scan. |
| `-workers` | `4` | Number of worker goroutines started to process files. |
| `-top` | `10` | Maximum number of distinct words retained **per file**. |

The program prints each file's base name, its total word count, and its retained word counts. It also prints elapsed time with `fmt.Println` and `time.Since`. Report order can change between runs because workers finish in different orders. The order of words inside a report can also change because the retained words are stored in a Go map.

There is currently no validation for empty `-dir`, nonpositive `-workers`, or nonpositive `-top`. Use a positive worker count and top value. Invalid values may hang or panic instead of producing a clear message.

## How the current workflow works

1. `run` reads the three CLI flags and calls `work`.
2. `readDirectory` uses `filepath.WalkDir` to visit the starting path and its subdirectories. It adds the discovered path of each regular file whose extension is exactly lowercase `.txt`. Other extensions, including `.TXT`, are skipped.
3. `work` creates a `jobs` channel of file paths and a `result` channel of `FileResult` values. It starts the requested number of workers, then a goroutine sends every discovered path into `jobs` and closes that channel.
4. Each worker receives paths from `jobs` until it closes. For each path, `process` opens the file, counts its words, keeps up to `-top` words for that file, and sends one `FileResult` into `result`. A file error produces an error result and the worker continues with later jobs.
5. A `sync.WaitGroup` tracks the workers. When all workers return, a separate goroutine closes `result`. Meanwhile, `work` receives results until the channel closes and returns the collected slice.
6. `run` prints each collected result. `main` prints the elapsed time and, if `run` returned an error, prints that error.

The `jobs` channel is buffered to twice the worker count; the `result` channel is unbuffered. The collector reads results while workers send them, so the workers can continue. Workers do **not** update a shared totals map, and the current code uses no mutex. There is no combined folder-wide word count.

## What each source file does

All seven files below are in `internal/analyzer` and use `package main`.

| File | Current responsibility |
| --- | --- |
| `flag.go` | Defines `serverConfig` and `loadConfig`. Registers `-dir`, `-workers`, and `-top`, sets usage text, parses flags, and returns the config. |
| `files.go` | `readDirectory` walks the directory tree and returns matching paths. `readFile` converts file bytes to text and calls `tokenize`. `fileName` takes the last slash-separated part of a path for display. |
| `words.go` | `tokenize` splits text wherever a rune is not a Unicode letter, then lowercases each word. Thus `Go, go!` becomes two `go` words, and `go123lang` becomes `go` and `lang`. |
| `stats.go` | Defines the per-file data types `FileStats` and `FileResult`, plus `countWord(io.Reader)`. Despite the file name, it does not keep shared folder statistics. |
| `pools.go` | Opens and processes files, runs the worker pool, sends jobs and results through channels, waits for workers, and collects results. |
| `sort.go` | Defines a temporary `wordCount` pair. `sortTop` sorts one file's distinct words by descending count and then alphabetically for ties, keeps at most `top`, and puts them back into a map. |
| `localMain.go` | Contains `run` and the program's `main`: invokes the workflow, prints per-file reports, and shows elapsed time or a returned error. |

`go.mod` defines the module as `GoProject1`. `internal/test` contains sample input files; files with other extensions are ignored by the analyzer. `internal/counter.go` is currently only a package declaration, not a working generic counter.

## Counting one file

`process` opens a path with `os.Open` and defers closing it. It passes the open file to `countWord`, whose `io.Reader` parameter means the same function can also accept another reader, such as a string reader in a future test. `countWord` currently uses `io.ReadAll` to load the **entire** file into memory. It passes the bytes through `readFile` and `tokenize`, increments a `map[string]int` for each word, and increments `total` for every occurrence. An empty file gives an empty map and a total of zero.

`FileStats` stores the display file name, that file's word-count map, and the total number of words before top-N filtering. `FileResult` stores a `FileStats` value plus an error. `sortTop` limits the map returned for display, but it does not change the total. If a file has fewer than `-top` distinct words, all its words are retained.

The temporary slice in `sortTop` has a predictable count-and-alphabetical order, but converting it back to a map loses that order. Consequently, the printed rows are **not** guaranteed to be ranked even though the set of retained top words is chosen by that ranking.

## Errors and current limitations

- If `filepath.WalkDir` fails, `readDirectory` wraps the error, `work` returns it, and `main` prints it. The program does not currently exit with a nonzero status for that error.
- If `os.Open` or `io.ReadAll` fails for one file, the worker sends an error result and continues. However, `run` currently ignores `FileResult.err` when printing, so a failed file appears as a blank-name, zero-total report. The error and failed path are not shown.
- A successful scan with no matching files produces no file reports; it prints only elapsed time.
- The scan accepts only regular files ending in lowercase `.txt`. It does not follow a symlink to a file as a regular file.
- `fileName` splits paths on `/` rather than using `filepath.Base`, which can affect display names on platforms with another path separator.
- There is no folder-wide count of files found, files processed, files failed, or total words; no merged word ranking; no tests in the current analyzer package; and no flag validation. The committed binary files are not needed for the `go run` command above.

## Project specification and remaining work

The original project goal is a CLI that scans a folder of text files, counts words, continues when an individual file fails, and prints a **folder-wide**, reproducible ranked report. The following describes that target behavior; see the sections above for what is already implemented.

- Discover regular `.txt` files recursively, including case-insensitive extensions. Report a missing or unreadable starting directory clearly. An empty folder should produce a valid zero-count report.
- Define a word as consecutive Unicode letters, lowercase it, and count every occurrence. For example, `Go, go! Café 123 café.` should yield `go=2`, `café=2`, and four total words. Keep the `io.Reader` interface so counting can be tested with strings as well as files.
- Process files through a configurable worker pool. Keep processing after a file fails, and include its path and error in the final report.
- Merge successful per-file counts into a folder-wide map and total. Report eligible files found, successful files, failed files, total words, and the requested number of most common words. Sort by descending count and alphabetically for ties, and print in that order. Results should be the same with one worker or several.
- Validate that `-dir` is not empty and `-workers` and `-top` are at least one. When fewer distinct words exist than `-top`, show only those available.
- Add meaningful tests for tokenization, empty files, filtering, errors, top-N bounds, and one-worker versus multiworker results. Run formatting, tests, and the race detector before calling the broader goal complete.

The original learning plan also includes shared-state protection with `sync.Mutex` if workers update shared totals, a reusable generic `Counter[T comparable]`, and a `FileStatus` enum pattern. Those are possible next steps; they are **not** present in the current implementation. The existing collector gathers separate results, so it does not currently need a mutex for its result slice.

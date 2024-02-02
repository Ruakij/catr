package main

import (
	"bufio"
	"bytes"
	"catr/textDetect"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/pflag"
)

var (
	listFlag           = pflag.BoolP("list", "l", false, "List files without displaying content")
	includeFlag        = pflag.StringArrayP("include", "i", []string{"*"}, "Include files")
	excludeFlag        = pflag.StringArrayP("exclude", "e", []string{""}, "Exclude files")
	formatFlag         = pflag.StringP("format", "o", "%s\n---\n%s\n---\n\n", "Customize how the file content is printed")
	textOnlyFlag       = pflag.Bool("text", true, "Only display text-files")
	ignoreEmptyFlag    = pflag.Bool("ignoreEmpty", true, "Ignore empty files")
	trimFileEnding     = pflag.Bool("trimFileEnding", true, "Trim newlines from end of files")
	parallelProcessing = pflag.Bool("parallel", false, "Parallel processing, faster for lots of files, but out-of-order")

	wg sync.WaitGroup
)

func main() {
	pflag.CommandLine.SortFlags = false
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nPrint path and content of all files recursively.\n\n %s [path ..]\n\n", filepath.Base(os.Args[0]))
		pflag.PrintDefaults()
	}

	pflag.Parse()

	paths := pflag.Args()
	if len(paths) < 1 {
		paths = []string{"."}
	}

	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "match glob %s: %s\n", path, err)
			continue
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				fmt.Fprintf(os.Stderr, "read %s: %s\n", path, err)
				continue
			}

			if info.IsDir() {
				err := walkAndMatch(match, *includeFlag, *excludeFlag)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			} else {
				printFileContent(match)
			}
		}
	}

	wg.Wait()
}

func walkAndMatch(inputLocation string, include []string, exclude []string) error {
	return filepath.WalkDir(inputLocation, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			//includeMatch := filepathMatchArray(include, path, entry, true, true)
			excludeMatch := filepathMatchArray(exclude, path, entry, true, true)
			//if includeMatch && !excludeMatch {
			if !excludeMatch {
				return nil
			}
			return fs.SkipDir
		}

		wg.Add(1)
		run := func() {
			defer wg.Done()

			includeMatch := filepathMatchArray(include, path, entry, true, true)
			excludeMatch := filepathMatchArray(exclude, path, entry, true, true)
			if includeMatch && !excludeMatch {
				printFileContent(path)
			}
		}
		if *parallelProcessing {
			go run()
		} else {
			run()
		}

		return nil
	})
}

func filepathMatchArray(matchStack []string, path string, entry fs.DirEntry, expected bool, anyMatch bool) bool {
	for _, match := range matchStack {
		// Suffix "/" depicts a directory-match
		if strings.HasSuffix(match, "/") {
			if entry.IsDir() {
				match = strings.TrimSuffix(match, "/")
			} else {
				continue
			}
		}

		// Prefix "/" depicts any match from the relative-root, skipping name-checks
		if strings.HasPrefix(match, "/") {
			match = strings.TrimPrefix(match, "/")
		} else {
			// Name-Match
			res2, _ := filepath.Match(match, entry.Name())
			if res2 == expected {
				if anyMatch {
					return true
				}
			} else if !anyMatch {
				return false
			}

			// path-matches against directories get a glob at the beginning if not already or has root-match
			if entry.IsDir() && !strings.HasPrefix(match, "/") && !strings.HasPrefix(match, "*") {
				match = "**/" + match
			}
		}

		// Path-Match
		res, _ := filepath.Match(match, path)
		if res == expected {
			if anyMatch {
				return true
			}
		} else if !anyMatch {
			return false
		}
	}
	return !anyMatch
}

func printFileContent(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	// Ignore empty files
	if fileInfo.Size() == 0 && *ignoreEmptyFlag {
		return
	}

	var content bytes.Buffer
	reader := bufio.NewReader(file)

	if *textOnlyFlag && fileInfo.Size() > 0 {
		buffer := make([]byte, 512)
		n, err := reader.Read(buffer)
		buffer = buffer[:n]
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %s\n", path, err)
			return
		}

		if textDetect.DetectEncoding(buffer) == textDetect.Unknown {
			return
		}

		content.Write(buffer)
	}

	if *listFlag {
		fmt.Println(path)
	} else {
		_, err = io.Copy(&content, reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %s\n", path, err)
		}

		// Strip the last newline
		if *trimFileEnding {
			contentBytes := content.Bytes()
			contentBytes = bytes.TrimRight(contentBytes, "\n ")
			content = *bytes.NewBuffer(contentBytes)
		}

		fmt.Printf(*formatFlag, path, content.String())
	}
}

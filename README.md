catr
--

Print path and content of all files recursively

<br>

<!-- TOC -->
- [1. Overview](#1-overview)
    - [1.1. Usage](#11-usage)
    - [1.2. Example](#12-example)
    - [1.3. Installation](#13-installation)
<!-- /TOC -->

<br>

# 1. Overview 

## 1.1. Usage

The program reads and prints the content of files based on specified options.

`catr <paths..>`

### 1.1.1. Options

| Option             | Description                               | Default        |
|--------------------|-------------------------------------------|----------------|
| `--list`, `-l`     | List files without displaying content     | `false`        |
| `--include`, `-i`  | Include filter                            | `["*"]`        |
| `--exclude`, `-e`  | Exclude filter                            | `[""]`         |
| `--format`, `-o`   | Customize how the file content is printed | `"%s\n---\n%s\n---\n\n"` |
| `--text`           | Only text-files (ASCII, UTF8, UTF16)      | `true`         |
| `--ignoreEmpty`    | Ignore empty files                        | `true`         |
| `--trimFileEnding` | Trim newlines from end of files           | `true`         |
| `--parallel`       | Parallel processing, but out-of-order     | `false`        |

<br>

## 1.2. Example

`catr *.go *.md`

`catr -l -i ".java" .`

<br>

## 1.3. Installation

To install the program, simply build it using the Go build tools and run the executable.

<br>


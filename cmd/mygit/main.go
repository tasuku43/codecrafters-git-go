package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		initGit()
	case "cat-file":
		catFile(os.Args[2:])
	case "hash-object":
		hashObject(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}

func initGit() {
	for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			os.Exit(1)
		}
	}

	headFileContents := []byte("ref: refs/heads/master\n")
	if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Initialized git directory")
}

func catFile(args []string) {
	switch option := args[0]; option {
	case "-p":
		hash := args[1]
		file, err := os.Open(fmt.Sprintf(".git/objects/%s/%s", hash[:2], hash[2:]))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()

		zr, err := zlib.NewReader(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating zlib reader: %s\n", err)
			os.Exit(1)
		}
		defer zr.Close()

		var content strings.Builder
		buf := make([]byte, 1024)
		for {
			n, err := zr.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
				os.Exit(1)
			}
			if n == 0 {
				break
			}
			content.Write(buf[:n])
		}

		splitContent := strings.SplitN(content.String(), "\x00", 2)
		if len(splitContent) != 2 {
			fmt.Fprintf(os.Stderr, "Error: Invalid object format\n")
			os.Exit(1)
		}
		fmt.Print(splitContent[1])
	default:
		fmt.Fprintf(os.Stderr, "Unknown option %s\n", option)
		os.Exit(1)
	}
}

func hashObject(args []string) {
	switch option := args[0]; option {
	case "-w":
		data, err := os.ReadFile(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
			os.Exit(1)
		}

		header := fmt.Sprintf("blob %d\x00", len(data))
		hasher := sha1.New()
		data = []byte(header + string(data))
		hasher.Write([]byte(data))
		hashBytes := hasher.Sum(nil)
		hesh := fmt.Sprintf("%x", hashBytes)

		var buffer bytes.Buffer
		w := zlib.NewWriter(&buffer)
		_, err = w.Write(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to zlib writer: %s\n", err)
			os.Exit(1)
		}
		err = w.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error closing zlib writer: %s\n", err)
			os.Exit(1)
		}

		filePath := fmt.Sprintf(".git/objects/%s/%s", hesh[:2], hesh[2:])
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(filePath, buffer.Bytes(), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "%s", hesh)
	default:
		fmt.Fprintf(os.Stderr, "Unknown option %s\n", option)
		os.Exit(1)
	}
}

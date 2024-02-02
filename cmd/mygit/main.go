package main

import (
	"fmt"
	"github.com/codecrafters-io/git-starter-go/pkg/gitcore"
	"os"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		gitcore.Init().Println()
	case "cat-file":
		gitcore.CatFile(os.Args[2:]).Print()
	case "hash-object":
		gitcore.HashObject(os.Args[2:]).Println()
	case "ls-tree":
		gitcore.LsTree(os.Args[2:]).Println()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}

package gitcore

import (
	"fmt"
	"os"
)

func Init() Message {
	directories := []string{".git", ".git/objects", ".git/refs"}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			handleError(err, fmt.Sprintf("Error creating directory: %s", dir))
		}
	}

	headFileContents := []byte("ref: refs/heads/master\n")
	if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
		handleError(err, "Error writing .git/HEAD file")
	}

	return Message("Initialized git directory")
}

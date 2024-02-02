package gitcore

import (
	"fmt"
	"os"
)

func handleError(err error, message string) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}

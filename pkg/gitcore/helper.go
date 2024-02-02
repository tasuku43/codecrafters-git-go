package gitcore

import (
	"fmt"
	"os"
)

type Message string
type Messages []Message

func (r Message) Print() {
	fmt.Print(r)
}

func (r Message) Println() {
	fmt.Println(r)
}

func (r Message) Bytes() []byte {
	return []byte(r)
}

func (r Messages) Print() {
	for _, message := range r {
		fmt.Print(message)
	}
}

func (r Messages) Println() {
	for _, message := range r {
		fmt.Println(message)
	}
}

func handleError(err error, message string) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}

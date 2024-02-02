package gitcore

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strings"
)

func CatFile(args []string) {
	if len(args) < 2 || args[0] != "-p" {
		handleError(fmt.Errorf("invalid usage"), "Usage: -p <hash>")
	}

	filePath := getObjectFilePath(args[1])
	contentBytes, err := decompressGitObject(filePath)
	if err != nil {
		handleError(err, "Error decompressing git object")
	}

	content := string(contentBytes)
	splitContent := strings.SplitN(content, "\x00", 2)
	if len(splitContent) != 2 {
		handleError(fmt.Errorf("invalid object format"), "")
	}

	fmt.Print(splitContent[1])
}

func getObjectFilePath(hash string) string {
	return fmt.Sprintf(".git/objects/%s/%s", hash[:2], hash[2:])
}

func decompressGitObject(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	zr, err := zlib.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	var result bytes.Buffer
	if _, err := io.Copy(&result, zr); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

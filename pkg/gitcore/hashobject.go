package gitcore

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
)

func HashObject(args []string) Message {
	if len(args) < 2 || args[0] != "-w" {
		handleError(fmt.Errorf("invalid usage"), "Usage: -w <file_path>")
	}

	data := readData(args[1])
	hash, fullData := calculateHash(data)
	compressedData := compressData(fullData)
	filePath := fmt.Sprintf(".git/objects/%s/%s", hash[:2], hash[2:])
	writeFile(filePath, compressedData)

	return Message(hash)
}

func readData(filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		handleError(err, "Error reading file")
	}
	return data
}

func calculateHash(data []byte) (string, []byte) {
	header := fmt.Sprintf("blob %d\x00", len(data))
	hasher := sha1.New()
	hasher.Write([]byte(header))
	hasher.Write(data)
	hashBytes := hasher.Sum(nil)
	hash := fmt.Sprintf("%x", hashBytes)
	return hash, append([]byte(header), data...)
}

func compressData(data []byte) []byte {
	var buffer bytes.Buffer
	w := zlib.NewWriter(&buffer)
	if _, err := w.Write(data); err != nil {
		handleError(err, "Error writing to zlib writer")
	}
	if err := w.Close(); err != nil {
		handleError(err, "Error closing zlib writer")
	}
	return buffer.Bytes()
}

func writeFile(filePath string, data []byte) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		handleError(err, "Error creating directory")
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		handleError(err, "Error writing file")
	}
}

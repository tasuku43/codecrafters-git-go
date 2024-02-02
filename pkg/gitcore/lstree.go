package gitcore

import (
	"bytes"
	"fmt"
)

type TreeEntry struct {
	mode []byte
	name []byte
	sha  []byte
}

func (t TreeEntry) string() string {
	return fmt.Sprintf("%s %s %x", t.mode, t.name, t.sha)
}

func LsTree(args []string) Messages {
	if len(args) < 2 || args[0] != "--name-only" {
		handleError(fmt.Errorf("invalid usage"), "Usage: --name-only <hash>")
	}
	treeBlob := CatFile([]string{"-p", args[1]}).Bytes()

	treeEntries := extractTreeEntries(treeBlob)

	var messages []Message
	for _, entry := range treeEntries {
		messages = append(messages, Message(entry.name))
	}
	return messages
}

func extractTreeEntries(treeBlob []byte) []TreeEntry {
	var entries []TreeEntry

	for len(treeBlob) > 0 {
		nullIndex := bytes.IndexByte(treeBlob, '\x00')
		if nullIndex == -1 || nullIndex+21 > len(treeBlob) {
			break
		}

		entryData := treeBlob[:nullIndex+21]
		entry := parseTreeEntry(entryData)
		entries = append(entries, entry)

		treeBlob = treeBlob[nullIndex+21:]
	}

	return entries
}

func parseTreeEntry(entryData []byte) TreeEntry {
	spaceIndex := bytes.IndexByte(entryData, ' ')
	nullIndex := bytes.IndexByte(entryData, '\x00')

	mode := entryData[:spaceIndex]
	name := entryData[spaceIndex+1 : nullIndex]
	sha := entryData[nullIndex+1 : nullIndex+21]

	return TreeEntry{mode: mode, name: name, sha: sha}
}

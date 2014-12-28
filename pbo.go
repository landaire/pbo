package pbo

import (
	"bufio"
	"os"
)

type Pbo struct {
	file            *os.File
	HeaderExtension *HeaderExtension
	Entries         []FileEntry
}

// Reads the file given by path and returns
// a Pbo pointer and err != nil if no errors occurred
func NewPbo(path string) (*Pbo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	pbo := Pbo{
		file: file,
	}

	// Create a new buffered reader
	reader := bufio.NewReader(file)

	for {
		entry := readEntry(reader)
		if entry.IsNull() {
			break
		} else if entry.Flag == ProductEntry {
			extension := HeaderExtension{
				FileEntry: entry,
			}

			extension.ReadExtendedFields(reader)
			pbo.HeaderExtension = &extension
		} else {
			pbo.Entries = append(pbo.Entries, entry)
		}

	}

	return &pbo, nil
}

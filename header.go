//go:generate stringer -type=Flag

package pbo

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"time"
)

type Flag uint32

const (
	Uncompressed Flag = 0x00000000
	Packed       Flag = 0x43707273
	ProductEntry Flag = 0x56657273 // resistance/elite/arma
)

type FileEntry struct {
	Name          string
	Flag          Flag
	OriginalSize  uint32
	ReservedField uint32
	Timestamp     time.Time
	DataSize      uint32
	data          *interface{}
}

func (f FileEntry) IsNull() bool {
	return f.Name == "" && f.Flag == 0 &&
		f.OriginalSize == 0 && f.ReservedField == 0 &&
		f.Timestamp.Unix() == 0 && f.DataSize == 0
}

func (f FileEntry) String() string {
	return fmt.Sprintf(
		"Name: %s\nFlag: 0x%08X (%s)\nOriginal Size: %d\nReserved: %d\nTimestamp: %d (%s)\nData Size: %d",
		f.Name,
		uint32(f.Flag), f.Flag.String(),
		f.OriginalSize,
		f.ReservedField,
		f.Timestamp.Unix(), f.Timestamp,
		f.DataSize,
	)
}

type HeaderExtension struct {
	FileEntry
	ExtendedFields map[string]string
}

func (e *HeaderExtension) ReadExtendedFields(r *bufio.Reader) {
	maybeNullByte, _ := r.Peek(1)
	if maybeNullByte[0] != 0 {
		e.ExtendedFields = make(map[string]string)
	}

	for maybeNullByte[0] != 0 {
		key, _ := r.ReadString(0)
		value, _ := r.ReadString(0)

		e.ExtendedFields[key[:len(key)-1]] = value[:len(value)-1]

		maybeNullByte, _ = r.Peek(1)
	}

	r.ReadByte()
}

func readEntry(r *bufio.Reader) FileEntry {
	entry := FileEntry{}
	entry.Name, _ = r.ReadString(0)
	entry.Name = entry.Name[:len(entry.Name)-1]

	var timestamp uint32
	fields := []interface{}{
		&entry.Flag,
		&entry.OriginalSize,
		&entry.ReservedField,
		&timestamp,
		&entry.DataSize,
	}

	for _, field := range fields {
		// Ignore errors -- swag
		binary.Read(r, binary.LittleEndian, field)
	}

	entry.Timestamp = time.Unix(int64(timestamp), 0)

	return entry
}

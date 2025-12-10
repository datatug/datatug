package dtio

import (
	"bytes"
	"os"
)

var sqliteHeader = []byte("SQLite format 3\x00")

func IsSQLite(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer func() {
		_ = f.Close()
	}()

	header := make([]byte, 16)
	n, err := f.Read(header)
	if err != nil || n < 16 {
		return false
	}

	return bytes.Equal(header, sqliteHeader)
}

package util

import (
	"os"
)

func WriteChunkToFile(file *os.File, data []byte) error {
	w := 0
	n := len(data)
	for {
		nw, err := file.Write(data[w:])
		if err != nil {
			return err
		}
		w += nw
		if nw >= n {
			return nil
		}
	}
}

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

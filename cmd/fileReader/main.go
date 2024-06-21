package fileReader

import (
	"io"
	"os"
)

func ReadFileContents(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

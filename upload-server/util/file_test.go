package util

import (
	"bufio"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileExists(t *testing.T) {
	t.Run("file exist", func(t *testing.T) {
		exist := FileExists("data/cute_cat.jpg")
		require.Equal(t, true, exist)
	})
	t.Run("file do not exist", func(t *testing.T) {
		exist := FileExists("data/cute.cat")
		require.Equal(t, false, exist)
	})
}

func TestWriteChunkToFile(t *testing.T) {
	t.Run("write chunk to file", func(t *testing.T) {
		// given os file
		file, err := os.Open("data/cute_cat.jpg")
		require.NoError(t, err)
		newFile, err := os.Create("data/cute_sad_cat.jpg")
		require.NoError(t, err)
		defer func() {
			err = file.Close()
			require.NoError(t, err)
			err = newFile.Close()
			require.NoError(t, err)
			err = os.Remove(newFile.Name())
			require.NoError(t, err)
		}()

		// when split file and copy chunks to another
		reader := bufio.NewReader(file)
		buf := make([]byte, 16)
		for {
			n, err := reader.Read(buf)
			if err != nil && err == io.EOF {
				break
			}
			err = WriteChunkToFile(newFile, buf[0:n])
			require.NoError(t, err)
		}

		// then
		fileInfo, err := file.Stat()
		require.NoError(t, err)
		newFileInfo, err := newFile.Stat()
		require.NoError(t, err)
		require.Equal(t, fileInfo.Size(), newFileInfo.Size())
	})
}

func TestGetFileContentType(t *testing.T) {
	t.Run("get file content Type", func(t *testing.T) {
		// given
		file, err := os.Open("data/cute_cat.jpg")
		require.NoError(t, err)
		defer file.Close()

		// when
		fileType, err := GetFileContentType(file)
		require.NoError(t, err)

		// then
		require.Equal(t, "image/jpeg", fileType)
	})

	t.Run("get file content Type pdf", func(t *testing.T) {
		// given
		file, err := os.Open("data/Ori_Coding_Test__1.0_.pdf")
		require.NoError(t, err)
		defer file.Close()

		// when
		fileType, err := GetFileContentType(file)
		require.NoError(t, err)

		// then
		require.Equal(t, "application/pdf", fileType)
	})

	t.Run("get file content Type unknown file", func(t *testing.T) {
		// given
		file, err := os.Open("data/crazy-file")
		require.NoError(t, err)
		defer file.Close()

		// when
		fileType, err := GetFileContentType(file)
		require.NoError(t, err)

		// then
		require.Equal(t, "application/octet-stream", fileType)
	})
}

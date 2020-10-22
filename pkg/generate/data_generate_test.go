package generate

import (
    "bytes"
    "github.com/bittygarden/lilac/pkg/io_tool"
    "github.com/stretchr/testify/assert"
    "os"
    "strconv"
    "testing"
)

func TestCreateFileIfNotExists(t *testing.T) {
    filePath := "d://test.txt"
    CreateFileIfNotExists(filePath)
    assert.True(t, io_tool.FileExists(filePath), "fail")
}

func TestAppendToFile(t *testing.T) {
    filePath := "d://test.txt"
    count := 10
    file, err := os.OpenFile(filePath, os.O_APPEND, os.ModeAppend)
    defer file.Close()

    assert.Nil(t, err, "打开文件异常", "err", err)
    for i := 0; i < count; i++ {
        buffer := bytes.Buffer{}
        buffer.WriteString(strconv.Itoa(i) + "\n")
        AppendToFile(&buffer, file)
    }
}

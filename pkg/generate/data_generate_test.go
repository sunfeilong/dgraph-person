package generate

import (
    "github.com/stretchr/testify/assert"
    "os"
    "strconv"
    "testing"
)

func TestCreateFileIfNotExists(t *testing.T) {
    filePath := "d://test.txt"
    CreateFileIfNotExists(filePath)
    assert.True(t, FileExist(filePath), "fail")
}

func TestAppendToFile(t *testing.T) {
    filePath := "d://test.txt"
    count := 10
    file, err := os.OpenFile(filePath, os.O_APPEND, os.ModeAppend)
    defer file.Close()

    assert.Nil(t, err, "打开文件异常", "err", err)
    for i := 0; i < count; i++ {
        AppendToFile(strconv.Itoa(i)+"\n", file)
    }
}

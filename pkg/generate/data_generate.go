package generate

import (
    "fmt"
    "github.com/xiaotian/dgraph-person/pkg/d_log"
    "github.com/xiaotian/dgraph-person/pkg/tools"
    "math/rand"
    "os"
    "strconv"
)

const APPEND_BARRIER int = 1024 * 1024 * 20 / 8

var logger = d_log.New()

func ProductData(totalCount int, maxFriendCount int, filePath string) {
    DeleteFileIfExists(filePath)
    CreateFileIfNotExists(filePath)

    file, err := os.OpenFile(filePath, os.O_APPEND, os.ModeAppend)
    if err != nil {
        logger.Errorw("GenerateData open file exception", "err", err)
    }
    defer file.Close()

    data := ""
    for i := 1; i <= totalCount; i++ {
        data = data + fmt.Sprintf("<%d> <num> \"%s\" .\n", i, strconv.Itoa(i))
        if needAdd(data) {
            AppendToFile(data, file)
            tools.ShowProgress("Generate person", i, totalCount)
            data = ""
        }
    }

    for i := 1; i <= totalCount; i++ {
        friends := getFriend(i, maxFriendCount, totalCount)
        for _, v := range *friends {
            data = data + fmt.Sprintf("<%d> <friend> <%d> .\n", i, v)
            data = data + fmt.Sprintf("<%d> <friend> <%d> .\n", v, i)
            if needAdd(data) {
                AppendToFile(data, file)
                tools.ShowProgress("Generate person friend", i, totalCount)
                data = ""
            }
        }
    }
    AppendToFile(data, file)
}

func getFriend(num int, maxFriendCount int, totalCount int) *[]int {
    if maxFriendCount <= 0 {
        return &[]int{}
    }
    rand.Seed(int64(num))
    friendCount := rand.Intn(maxFriendCount + 1)
    result := make([]int, friendCount)
    for i := 0; i < friendCount; i++ {
        result[i] = rand.Intn(totalCount + 1)
    }
    return &result
}

func needAdd(data string) bool {
    return len(data) >= APPEND_BARRIER
}

func AppendToFile(data string, file *os.File) {
    n, err := file.Write([]byte(data))
    if err != nil {
        logger.Errorw("append to file Exception", "err", err)
    }
    logger.Infow("append to file.", "length", n)
}

func FileExist(path string) bool {
    _, err := os.Lstat(path)
    return !os.IsNotExist(err)
}

func CreateFileIfNotExists(filePath string) {
    if !FileExist(filePath) {
        logger.Infow("file not exists, create file.", "file", filePath)
        file, err := os.Create(filePath)
        if err != nil {
            logger.Errorw("create file exception", "err", err)
        }
        defer file.Close()
    }
}
func DeleteFileIfExists(filePath string) {
    if FileExist(filePath) {
        logger.Infow("file exists, delete file.", "file", filePath)
        err := os.Remove(filePath)
        if err != nil {
            logger.Errorw("delete file exception", "err", err)
        }
    }
}

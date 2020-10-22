package generate

import (
    "bytes"
    "fmt"
    "github.com/bittygarden/lilac/pkg/io_tool"
    "github.com/xiaotian/dgraph-person/pkg/d_log"
    "github.com/xiaotian/dgraph-person/pkg/tools"
    "math/rand"
    "os"
    "strconv"
)

const APPEND_BARRIER int = 1024 * 1024 * 20 / 8

var logger = d_log.New()

var friendsMap = make(map[int]bool)

func ProductData(totalCount int, maxFriendCount int, filePath string) {
    DeleteFileIfExists(filePath)
    CreateFileIfNotExists(filePath)

    file, err := os.OpenFile(filePath, os.O_APPEND, os.ModeAppend)
    if err != nil {
        logger.Errorw("GenerateData open file exception", "err", err)
    }
    defer file.Close()

    data := bytes.Buffer{}
    for i := 1; i <= totalCount; i++ {
        data.WriteString(fmt.Sprintf("<%d> <num> \"%s\" .\n", i, strconv.Itoa(i)))
        if needAdd(&data) {
            AppendToFile(&data, file)
            tools.ShowProgress("Generate person", i, totalCount)
            data.Reset()
        }
    }

    for i := 1; i <= totalCount; i++ {
        friends := getFriend(i, maxFriendCount, totalCount)
        for _, v := range *friends {
            data.WriteString(fmt.Sprintf("<%d> <friend> <%d> .\n", i, v))
            data.WriteString(fmt.Sprintf("<%d> <friend> <%d> .\n", v, i))
            if needAdd(&data) {
                AppendToFile(&data, file)
                tools.ShowProgress("Generate person friend", i, totalCount)
                data.Reset()
            }
        }
    }
    AppendToFile(&data, file)
}

func getFriend(num int, maxFriendCount int, totalCount int) *[]int {
    if maxFriendCount <= 0 {
        return &[]int{}
    }
    rand.Seed(int64(num))
    friendCount := rand.Intn(maxFriendCount) + 1
    result := make([]int, friendCount)
    for ; ; {
        friend := 0
        index := 0
        for ; ; {
            friend = rand.Intn(totalCount) + 1
            if friendsMap[friend] {
                continue
            }
            result[index] = friend
            index++
            friendsMap[friend] = true
            if index == friendCount {
                for _, v := range result {
                    friendsMap[v] = false
                }
                return &result
            }
        }
    }
}

func needAdd(data *bytes.Buffer) bool {
    return data.Len() >= APPEND_BARRIER
}

func AppendToFile(data *bytes.Buffer, file *os.File) {
    n, err := file.Write(data.Bytes())
    if err != nil {
        logger.Errorw("append to file Exception", "err", err)
    }
    logger.Infow("append to file.", "length", n)
}

func CreateFileIfNotExists(filePath string) {
    if !io_tool.FileExists(filePath) {
        logger.Infow("file not exists, create file.", "file", filePath)
        file, err := os.Create(filePath)
        if err != nil {
            logger.Errorw("create file exception", "err", err)
        }
        defer file.Close()
    }
}
func DeleteFileIfExists(filePath string) {
    if io_tool.FileExists(filePath) {
        logger.Infow("file exists, delete file.", "file", filePath)
        err := os.Remove(filePath)
        if err != nil {
            logger.Errorw("delete file exception", "err", err)
        }
    }
}

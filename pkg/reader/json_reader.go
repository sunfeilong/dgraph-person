package reader

import (
    "encoding/json"
    "github.com/xiaotian/dgraph-person/pkg/d_log"
    "io/ioutil"
)

var logger = d_log.New()

type JsonPersonReader struct {
}

type IdNameAndPhone struct {
    Id    string `json:"id"`
    Name  string `json:"name"`
    Phone string `json:"phone"`
}

func (read JsonPersonReader) ReadFromFile(filePath string) []IdNameAndPhone {
    logger.Infow("JsonPersonReader ReadFromFile ", "filePath", filePath)

    dataBytes, err := ioutil.ReadFile(filePath)
    if err != nil {
        logger.Panicw("JsonPersonReader ReadFromFile error", "filePath", filePath, "err", err)
    }

    result := []IdNameAndPhone{}
    err = json.Unmarshal(dataBytes, &result)
    if err != nil {
        logger.Panicw("JsonPersonReader ReadFromFile Unmarshal error", "filePath", filePath, "err", err)
    }
    logger.Infof("JsonPersonReader ReadFromFile data length: %s", len(result))
    return result
}

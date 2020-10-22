package reader

import (
    "encoding/json"
    "github.com/xiaotian/dgraph-person/pkg/d_log"
    "io/ioutil"
)

var logger = d_log.New()

type JsonPersonReader struct {
}

type Num struct {
    Num  string `json:"num"`
}

type NumAndFriendNum struct {
    UserNum string `json:"userNum"`
    FriendNum     string `json:"friendNum"`
}

func (read JsonPersonReader) ReadPersonFromFile(filePath string) []Num {
    logger.Infow("JsonPersonReader ReadPersonFromFile ", "filePath", filePath)

    dataBytes, err := ioutil.ReadFile(filePath)
    if err != nil {
        logger.Panicw("JsonPersonReader ReadPersonFromFile error", "filePath", filePath, "err", err)
    }

    result := []Num{}
    err = json.Unmarshal(dataBytes, &result)
    if err != nil {
        logger.Panicw("JsonPersonReader ReadPersonFromFile Unmarshal error", "filePath", filePath, "err", err)
    }
    logger.Infof("JsonPersonReader ReadPersonFromFile data length: %s", len(result))
    return result
}

func (read JsonPersonReader) ReadFriendFromFile(filePath string) []NumAndFriendNum {
    logger.Infow("JsonPersonReader ReadFriendFromFile ", "filePath", filePath)

    dataBytes, err := ioutil.ReadFile(filePath)
    if err != nil {
        logger.Panicw("JsonPersonReader ReadFriendFromFile error", "filePath", filePath, "err", err)
    }

    result := []NumAndFriendNum{}
    err = json.Unmarshal(dataBytes, &result)
    if err != nil {
        logger.Panicw("JsonPersonReader ReadFriendFromFile Unmarshal error", "filePath", filePath, "err", err)
    }
    logger.Infof("JsonPersonReader ReadFriendFromFile data length: %s", len(result))
    return result
}

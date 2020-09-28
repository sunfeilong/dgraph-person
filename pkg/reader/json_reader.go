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
    Name  string `json:"name"`
    Phone string `json:"phone"`
}

type PhoneNameAndPhone struct {
    UserPhone string `json:"userPhone"`
    Name      string `json:"name"`
    Phone     string `json:"phone"`
}

func (read JsonPersonReader) ReadPersonFromFile(filePath string) []IdNameAndPhone {
    logger.Infow("JsonPersonReader ReadPersonFromFile ", "filePath", filePath)

    dataBytes, err := ioutil.ReadFile(filePath)
    if err != nil {
        logger.Panicw("JsonPersonReader ReadPersonFromFile error", "filePath", filePath, "err", err)
    }

    result := []IdNameAndPhone{}
    err = json.Unmarshal(dataBytes, &result)
    if err != nil {
        logger.Panicw("JsonPersonReader ReadPersonFromFile Unmarshal error", "filePath", filePath, "err", err)
    }
    logger.Infof("JsonPersonReader ReadPersonFromFile data length: %s", len(result))
    return result
}

func (read JsonPersonReader) ReadFriendFromFile(filePath string) []PhoneNameAndPhone {
    logger.Infow("JsonPersonReader ReadFriendFromFile ", "filePath", filePath)

    dataBytes, err := ioutil.ReadFile(filePath)
    if err != nil {
        logger.Panicw("JsonPersonReader ReadFriendFromFile error", "filePath", filePath, "err", err)
    }

    result := []PhoneNameAndPhone{}
    err = json.Unmarshal(dataBytes, &result)
    if err != nil {
        logger.Panicw("JsonPersonReader ReadFriendFromFile Unmarshal error", "filePath", filePath, "err", err)
    }
    logger.Infof("JsonPersonReader ReadFriendFromFile data length: %s", len(result))
    return result
}

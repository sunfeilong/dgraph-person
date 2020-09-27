package main

import (
    "github.com/xiaotian/dgraph-person/pkg/client"
    "github.com/xiaotian/dgraph-person/pkg/d_log"
    "github.com/xiaotian/dgraph-person/pkg/data"
    "regexp"
    "strings"
)

const pattern string = "\\d+"

var logger = d_log.New()

func main() {
    logger.Info("DGraph Person server start")
    c := client.Client{}
    c.Connect("10.0.8.36", 19080)
    initSchema(c)
    addPerson(c, "公安局", "110")
    person := c.GetPersonByEdge("phone", "15210012050")
    c.DeleteByUid(person.Uid)
    logger.Info("DGraph Person server end")
}

func initSchema(c client.Client) {
    _, err := c.AddSchema(data.Schema)
    if err != nil {
        panic(err)
    }
}

func addPerson(c client.Client, name string, phone string) {
    if len(strings.TrimSpace(name)) == 0 {
        return
    }
    ok, err := regexp.Match(pattern, []byte(phone))
    if err != nil && !ok {
        return
    }

    person := c.GetPersonByEdge("phone", phone)
    if person != nil {
        return
    }
    b, err := c.AddPerson(name, phone)
    if err != nil {
        panic(err)
    }
    if !b {
        panic("添加数据出错")
    }
}

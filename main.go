package main

import (
    "github.com/xiaotian/dgraph-person/pkg/client"
    "github.com/xiaotian/dgraph-person/pkg/data"
    "regexp"
    "strings"
)

const pattern string = "\\d+"

type Friend struct {
    FriendUserId string         `json:"userId"`
    Friends      []NameAndPhone `json:"friends"`
}

type NameAndPhone struct {
    Name  string `json:"name"`
    Phone string `json:"phone"`
}

type User struct {
    Uid   string   `json:"uid,omitempty"`
    Name  string   `json:"name,omitempty"`
    Phone string   `json:"phone,omitempty"`
    DType []string `json:"dgraph.type,omitempty"`
}

func main() {

    client := client.Client{}
    client.Connect("10.0.8.36", 19080)
    initSchema(client)
    //addPerson(client, "孙飞龙", "15210012054")
    person := client.GetPersonByEdge("phone", "15210012050")
    client.DeleteByUid(person.Uid)
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

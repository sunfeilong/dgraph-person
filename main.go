package main

import (
    "github.com/xiaotian/dgraph-person/pkg/client"
    "github.com/xiaotian/dgraph-person/pkg/d_log"
    "github.com/xiaotian/dgraph-person/pkg/data"
    "github.com/xiaotian/dgraph-person/pkg/reader"
    "github.com/xiaotian/dgraph-person/pkg/tools"
    "strings"
)

var logger = d_log.New()

func main() {
    userFile := "d://user.json"
    ip := "10.0.8.36"
    port := 19080

    logger.Info("DGraph Person server start", "ip", ip, "port", port, "userFile", userFile)
    c := client.Client{}
    c.Connect(ip, port)

    logger.Infow("init schema")
    initSchema(c)

    logger.Infow("read user data from file ", "file", userFile)
    jsonPersonReader := reader.JsonPersonReader{}
    idNameAndPhoneList := jsonPersonReader.ReadFromFile(userFile)
    length := len(idNameAndPhoneList)
    logger.Infow("read user data from file ", "file", userFile, "length", length)
    deletePersonList(c, getPhoneList(idNameAndPhoneList))

    //for i, v := range idNameAndPhoneList {
    //    tools.ShowProgress("AddPersonToData", i+1, length)
    //
    //    deletePerson(c, v.Phone)
    //}
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
        logger.Errorw("addPerson name is blank", "name", name)
        return
    }
    if !tools.IsNumber(phone) {
        logger.Panicw("addPerson phone is not a number", "phone", phone)
        return
    }

    person := c.GetPersonByEdge("phone", phone)
    if person != nil {
        logger.Infow("addPerson Person Has Exists ", "name", name, "phone", phone)
        return
    }
    b, err := c.AddPerson(name, phone)
    if err != nil {
        logger.Panicw("addPerson error", err, err)
    }
    if !b {
        logger.Panicw("addPerson error", err, err)
    }
}

func deletePerson(c client.Client, phone string) {
    if len(strings.TrimSpace(phone)) == 0 {
        logger.Errorw("deletePerson phone is blank", "phone", phone)
        return
    }

    person := c.GetPersonByEdge("phone", phone)
    if person == nil {
        logger.Infow("deletePerson Person Not Exists ", "phone", phone)
        return
    }
    b := c.DeleteByUid([]string{person.Uid})
    if !b {
        logger.Panicw("deletePerson error")
    }
}

func deletePersonList(c client.Client, phoneList []string) {
    if len(phoneList) == 0 {
        logger.Errorw("deletePerson phoneList is blank", "phoneList", phoneList)
        return
    }

    uidList := make([]string, len(phoneList))
    for i, v := range phoneList {
        person := c.GetPersonByEdge("phone", v)
        if person == nil {
            continue
        }
        uidList[i] = person.Uid
    }

    b := c.DeleteByUid(uidList)
    if !b {
        logger.Panicw("deletePerson error")
    }
}

func getPhoneList(idNamePhone []reader.IdNameAndPhone) []string {
    if len(idNamePhone) == 0 {
        return make([]string, 0)
    }

    phones := make([]string, len(idNamePhone))

    for i, v := range idNamePhone {
        phones[i] = v.Phone
    }

    return phones
}

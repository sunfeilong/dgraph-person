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
    friendFile := "d://friends.json"
    ip := "10.0.8.36"
    port := 19080

    logger.Info("DGraph Person server start", "ip", ip, "port", port, "userFile", userFile, "friendFile", friendFile)
    c := client.Client{}
    c.Connect(ip, port)
    logger.Infow("init schema")
    updateSchema(c)
    //loadPerson(c, userFile)
    loadFriend(c, friendFile)

    logger.Info("DGraph Person server end")
}

func loadFriend(c client.Client, friendFile string) {
    logger.Infow("loadFriend load data file ", "file", friendFile)
    jsonPersonReader := reader.JsonPersonReader{}
    phoneNameAndPhoneList := jsonPersonReader.ReadFriendFromFile(friendFile)
    length := len(phoneNameAndPhoneList)
    logger.Infow("loadFriend load data file ", "file", friendFile, "length", length)
    for i, v := range phoneNameAndPhoneList {
        if !addPerson(c, v.Name, v.Phone) {
            logger.Infow("loadFriend add person to db failed ", "name", v.Name, "phone", v.Phone)
            continue
        }
        friend := c.GetPersonByEdge("phone", v.Phone)
        if nil == friend {
            logger.Infow("loadFriend get friend from db, friend not exists ", "file", friendFile, "length", length)
            continue
        }
        person := c.GetPersonByEdge("phone", v.UserPhone)
        if nil == person {
            logger.Infow("loadFriend from file get person from db, person not exists ", "file", friendFile, "length", length)
            continue
        }
        //add relation
        c.AddFriend(friend.Uid, person.Uid)
        tools.ShowProgress("LoadFriend", i+1, length)
    }
}

func loadPerson(c client.Client, userFile string) {
    logger.Infow("read user data from file ", "file", userFile)
    jsonPersonReader := reader.JsonPersonReader{}
    idNameAndPhoneList := jsonPersonReader.ReadPersonFromFile(userFile)
    length := len(idNameAndPhoneList)
    logger.Infow("read user data from file ", "file", userFile, "length", length)
    for i, v := range idNameAndPhoneList {
        addPerson(c, v.Name, v.Phone)
        tools.ShowProgress("AddPersonToData", i+1, length)
    }
}

func updateSchema(c client.Client) {
    _, err := c.AddSchema(data.Schema)
    if err != nil {
        panic(err)
    }
}

func addPerson(c client.Client, name string, phone string) bool {
    if len(strings.TrimSpace(name)) == 0 {
        logger.Errorw("addPerson name is blank", "name", name)
        return false
    }
    if !tools.IsNumber(phone) {
        logger.Panicw("addPerson phone is not a number", "phone", phone)
        return false
    }

    person := c.GetPersonByEdge("phone", phone)
    if person != nil {
        logger.Infow("addPerson Person Has Exists ", "name", name, "phone", phone)
        return true
    }
    b, err := c.AddPerson(name, phone)
    if err != nil {
        logger.Panicw("addPerson error", err, err)
    }
    if !b {
        logger.Panicw("addPerson error", err, err)
    }
    return true
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

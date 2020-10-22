package main

import (
    "flag"
    "fmt"
    "github.com/xiaotian/dgraph-person/pkg/client"
    "github.com/xiaotian/dgraph-person/pkg/d_log"
    "github.com/xiaotian/dgraph-person/pkg/data"
    "github.com/xiaotian/dgraph-person/pkg/generate"
    "github.com/xiaotian/dgraph-person/pkg/reader"
    "github.com/xiaotian/dgraph-person/pkg/tools"
    "strings"
)

var logger = d_log.New()

func main() {

    command := flag.String("command", "", "run command \n[generate] generate data\n[clean] clean all data in dgraph\n[update] update schema")
    host := flag.String("host", "", "[clean/update] dgraph host")
    port := flag.Int("port", 0, "[clean/update] dgraph port")
    filePath := flag.String("filePath", "", "[generate] save generate data file path")
    count := flag.Int("count", 1000000, "[generate] person count")
    maxFriendCount := flag.Int("maxFriendCount", 200, "[generate] person max friend count")

    flag.Parse()

    if command == nil || *command == "" {
        fmt.Println("command must in [generate clean update]")
    } else if *command == "generate" {
        if checkIsEmpty("count", host) {
            return
        }
        if *count <= 0 {
            fmt.Println("count must grater than 0", "count", *count)
        }
        if *maxFriendCount <= 0 {
            fmt.Println("count must grater than 0", "maxFriendCount", *maxFriendCount)
        }
        if * maxFriendCount >= *count {
            fmt.Println("maxFriendCount must less than count", "count", *count, "maxFriendCount", maxFriendCount)
        }
        if checkIsEmpty("maxFriendCount", port) {
            return
        }
        if checkIsEmpty("filePath", filePath) {
            return
        }
        generateData(*count, *maxFriendCount, *filePath)
    } else if *command == "clean" {
        fmt.Printf("clean all dgraph data start. hosy:%s, port:%d\n", *host, *port)
        if checkIsEmpty("host", host) {
            return
        }
        if checkIsEmpty("port", port) {
            return
        }
        if *host == "" || *port == 0 {
            fmt.Printf("host and port must exists. host:%s, port:%d\n", *host, *port)
            return
        }
        c := client.Client{}
        c.Connect(*host, *port)
        dropAll(c)
        fmt.Printf("clean all dgraph data end. host:%s, port:%d\n", *host, *port)
    } else if *command == "update" {
        fmt.Println("update schema start.")
        if checkIsEmpty("host", host) {
            return
        }
        if checkIsEmpty("port", port) {
            return
        }
        c := client.Client{}
        c.Connect(*host, *port)
        updateSchema(c)
        fmt.Println("update schema success.")
    } else {
        fmt.Println("command must in [generate/clean/update]", "\nyou input command: ", *command)
    }
}

func checkIsEmpty(name string, value interface{}) bool {
    if value == nil {
        fmt.Printf("%s must not nil.\n", name)
        return true
    }
    return false
}

func generateData(totalCount int, maxFriendCount int, filePath string) {
    generate.ProductData(totalCount, maxFriendCount, filePath)
}

func dropAll(c client.Client) {
    logger.Infow("drop all data")
    c.DropAll()
}

func loadFriend(c client.Client, friendFile string) {
    logger.Infow("loadFriend load data file ", "file", friendFile)
    jsonPersonReader := reader.JsonPersonReader{}
    phoneNameAndPhoneList := jsonPersonReader.ReadFriendFromFile(friendFile)
    length := len(phoneNameAndPhoneList)
    logger.Infow("loadFriend load data file ", "file", friendFile, "length", length)
    for i, v := range phoneNameAndPhoneList {
        logger.Infow("loadFriend friend", "num", v.UserNum, "friendNum", v.FriendNum)
        if !addPerson(c, v.FriendNum) {
            logger.Infow("loadFriend add person to db failed ", "num", v.UserNum, "friendNum", v.FriendNum)
            continue
        }
        friend := c.GetPersonByEdge("num", v.FriendNum)
        if nil == friend {
            logger.Infow("loadFriend get friend from db, friend not exists ", "file", friendFile, "length", length)
            continue
        }
        person := c.GetPersonByEdge("num", v.UserNum)
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
        addPerson(c, v.Num)
        tools.ShowProgress("AddPersonToData", i+1, length)
    }
}

func updateSchema(c client.Client) {
    logger.Infow("init schema")
    _, err := c.AddSchema(data.Schema)
    if err != nil {
        panic(err)
    }
}

func addPerson(c client.Client, num string) bool {
    if len(strings.TrimSpace(num)) == 0 {
        logger.Errorw("addPerson name is blank", "num", num)
        return false
    }
    if !tools.IsNumber(num) {
        logger.Errorw("addPerson num is not a number", "num", num)
        return false
    }

    person := c.GetPersonByEdge("num", num)
    if person != nil {
        logger.Infow("addPerson Person Has Exists ", "num", num)
        return true
    }
    b, err := c.AddPerson(num)
    if err != nil {
        logger.Panicw("addPerson error", err, err)
    }
    if !b {
        logger.Panicw("addPerson error", err, err)
    }
    return true
}

func deletePerson(c client.Client, num string) {
    if len(strings.TrimSpace(num)) == 0 {
        logger.Errorw("deletePerson num is blank", "num", num)
        return
    }

    person := c.GetPersonByEdge("num", num)
    if person == nil {
        logger.Infow("deletePerson Person Not Exists ", "num", num)
        return
    }
    b := c.DeleteByUid([]string{person.Uid})
    if !b {
        logger.Panicw("deletePerson error")
    }
}

func deletePersonList(c client.Client, numList []string) {
    if len(numList) == 0 {
        logger.Errorw("deletePerson numList is blank", "numList", numList)
        return
    }

    uidList := make([]string, len(numList))
    for i, v := range numList {
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

func getNumList(numList []reader.Num) []string {
    if len(numList) == 0 {
        return make([]string, 0)
    }

    nums := make([]string, len(numList))

    for i, v := range numList {
        nums[i] = v.Num
    }
    return nums
}

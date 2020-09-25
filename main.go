package main

import (
	"context"
	"dgraph/pkg/client"
	"dgraph/pkg/data"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"strings"
)

const queryTemplate = `
		{
			user(func: eq(phone, "_phone")) {
                uid
				name
				phone
				dgraph.type
			}
		}
	`

type Records struct {
	Records []Person `json:"RECORDS"`
}

type Person struct {
	UserId string `json:"user_id"`
	Name   string `json:"real_name"`
	Phone  string `json:"phone"`
}

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
	//userDataFile := "d://user.json"
	//
	//friendsDataFile := "d://friends.json"
	//
	//userData, err := ioutil.ReadFile(userDataFile)
	//if nil != err {
	//	panic(err)
	//}
	//
	//friendData, err := ioutil.ReadFile(friendsDataFile)
	//if nil != err {
	//	panic(err)
	//}
	//
	//records := &Records{}
	//err = json.Unmarshal(userData, records)
	//
	//friends := &[]Friend{}
	//err = json.Unmarshal(friendData, friends)
	//
	//conn, err := grpc.Dial("10.0.8.36:19080", grpc.WithInsecure())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer conn.Close()
	//
	//dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	//
	//for _, d := range records.Records {
	//	deletePhone(dg, d.Phone)
	//}
	//
	//for _, f := range *friends {
	//	if nil != f.Friends {
	//		for _, p := range f.Friends {
	//			deletePhone(dg, p.Phone)
	//		}
	//	}
	//
	//}

	client := client.Client{}
	client.Connect("10.0.8.36", 19080)

	schema, err := client.AddSchema(data.Schema)
	if err != nil {
		panic(err)
	}
	fmt.Println(schema)

}

func checkAndAdd(c *dgo.Dgraph, name string, phone string) {
	if len(strings.TrimSpace(name)) == 0 {
		fmt.Println("name is empty, phone: ", phone)
		return
	}

	data := getData(c, phone)
	u := User{
		Uid:   "_:u",
		Name:  name,
		Phone: phone,
		DType: []string{"User"},
	}

	if data != nil && len(strings.TrimSpace(data.Name)) != 0 {
		fmt.Println("data has exists and name is not empty, phone: ", phone, "name: ", name)
		return
	}

	if data != nil {
		fmt.Println("user has exists, phone: " + phone)
		u.Uid = data.Uid
	}

	da, err := json.Marshal(u)

	if nil != err {
		panic(err)
	}

	fmt.Println(string(da))
	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   da,
	}

	res, err := c.NewTxn().Mutate(context.Background(), mu)
	if err != nil {
		panic(err)
	}

	fmt.Println(&res.Json)
}

func addRelation(c *dgo.Dgraph, phoneA string, phoneB string) {

}

func deletePhone(c *dgo.Dgraph, phone string) {
	data := getData(c, phone)
	if data == nil {
		return
	}

	d := map[string]string{"uid": data.Uid}
	marshal, err := json.Marshal(d)
	if nil != err {

		panic(err)
	}

	mu := &api.Mutation{
		CommitNow:  true,
		DeleteJson: marshal,
	}

	mutate, err := c.NewTxn().Mutate(context.Background(), mu)
	if nil != err {
		panic(err)
	}

	fmt.Println(mutate.Json)

}

func getData(c *dgo.Dgraph, phone string) *User {
	query := strings.Replace(queryTemplate, "_phone", phone, 1)
	txn := c.NewTxn()

	resp, err := txn.QueryWithVars(context.Background(), query, nil)
	if err != nil {
		panic(err)
	}

	var decode struct {
		User []User `json:"user,omitempty"`
	}

	err = json.Unmarshal(resp.GetJson(), &decode)
	fmt.Println(string(resp.Json))
	if err != nil {
		panic(err)
	}

	if len(decode.User) > 1 {
		panic("")
	}

	if len(decode.User) == 1 {
		return &decode.User[0]
	}

	return nil
}

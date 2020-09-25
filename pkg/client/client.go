package client

import (
    "context"
    "encoding/json"
    "github.com/dgraph-io/dgo/v200"
    "github.com/dgraph-io/dgo/v200/protos/api"
    "github.com/xiaotian/dgraph-person/pkg/data"
    "google.golang.org/grpc"
    "strings"
)

type Client struct {
    ip   string
    port int
    dg   *dgo.Dgraph
}

type DecodePerson struct {
    P []data.Person `json:"data,omitempty"`
}

func (c *Client) Connect(ip string, port int) {
    c.ip = ip
    c.port = port
    dial, err := grpc.Dial("10.0.8.36:19080", grpc.WithInsecure())
    if err != nil {
        panic(err)
    }
    c.dg = dgo.NewDgraphClient(api.NewDgraphClient(dial))
}

func (c *Client) AddSchema(schema string) (bool, error) {
    op := &api.Operation{
        Schema:          schema,
        RunInBackground: false,
    }
    err := c.dg.Alter(context.Background(), op)
    if err != nil {
        return false, err
    }
    return true, nil
}

func (c *Client) AddPerson(name string, phone string) (bool, error) {

    addData := data.Person{
        Uid:      "_:person",
        Name:     name,
        Phone:    phone,
        DType:    []string{"Person"},
        NodeType: "Person",
    }

    jsonData, err := json.Marshal(addData)
    if nil != err {
        panic(err)
    }
    mu := &api.Mutation{
        CommitNow: true,
        SetJson:   jsonData,
    }

    response, err := c.dg.NewTxn().Mutate(context.Background(), mu)
    if err != nil {
        return false, err
    }
    if response.Json != nil {
        return false, nil
    }
    return true, nil
}

func getByUid(uid string, result *interface{}) {

}

func (c *Client) GetPersonByEdge(edgeName string, value string) *data.Person {
    txn := c.dg.NewTxn()
    ctx := context.Background()
    q := strings.Replace(strings.Replace(data.QueryByEdge, "$edgeName", edgeName, 1), "$phone", value, 1)
    res, err := txn.Query(ctx, q)
    if err != nil {
        panic(err)
    }
    decode := DecodePerson{}
    err = json.Unmarshal(res.Json, &decode)
    if len(decode.P) == 0 {
        return nil
    } else if len(decode.P) == 1 {
        return &decode.P[0]
    } else {
        panic("数据存在错误")
    }
}

func deleteByUid(uid string) bool {
    return false
}

func deleteEdge(edgeName string, value string, result *interface{}) {

}

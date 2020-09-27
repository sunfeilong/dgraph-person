package client

import (
    "context"
    "encoding/json"
    "github.com/dgraph-io/dgo/v200"
    "github.com/dgraph-io/dgo/v200/protos/api"
    "github.com/xiaotian/dgraph-person/pkg/d_log"
    "github.com/xiaotian/dgraph-person/pkg/data"
    "google.golang.org/grpc"
    "strings"
)

var logger = d_log.New()

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
    logger.Infow("connect to DGraph server start", "ip", ip, "port", port)
    dial, err := grpc.Dial(ip+":"+string(port), grpc.WithInsecure())
    if err != nil {
        logger.Panicw("connect to DGraph server error", "ip", ip, "port", port, "err", err.Error())
    }
    c.dg = dgo.NewDgraphClient(api.NewDgraphClient(dial))
    logger.Infow("connect to DGraph server success", "ip", ip, "port", port)
}

func (c *Client) AddSchema(schema string) (bool, error) {
    op := &api.Operation{
        Schema:          schema,
        RunInBackground: false,
    }
    logger.Infow("AddSchema start", "schema", schema)
    err := c.dg.Alter(context.Background(), op)
    if err != nil {
        logger.Errorw("AddSchema error", "schema", schema, "err", err.Error())
        return false, err
    }
    logger.Infow("AddSchema start success", "schema", schema)
    return true, nil
}

func (c *Client) AddPerson(name string, phone string) (bool, error) {

    logger.Infow("AddPerson start", "name", name, "phone", phone)

    addData := data.Person{
        Uid:      "_:person",
        Name:     name,
        Phone:    phone,
        DType:    []string{"Person"},
        NodeType: "Person",
    }
    jsonData, err := json.Marshal(addData)
    if nil != err {
        logger.Errorw("AddPerson json marshall error", "err", err)
        return false, err
    }

    mu := &api.Mutation{
        CommitNow: true,
        SetJson:   jsonData,
    }

    ctx := context.Background()
    txn := c.dg.NewTxn()
    response, err := txn.Mutate(ctx, mu)
    defer txn.Discard(ctx)

    if err != nil {
        logger.Errorw("AddPerson mutation error", "name", name, "phone", phone, "err", err)
        return false, err
    }

    logger.Infow("AddPerson mutation response", "name", name, "phone", phone, "response", string(response.Json))
    if response.Json != nil {
        return false, nil
    }
    logger.Infow("AddPerson mutation success", "name", name, "phone", phone, "response", string(response.Json))
    return true, nil
}

func (c *Client) GetPersonByUid(uid string) *data.Person {
    logger.Infow("GetPersonByUid start", "uid", uid)

    txn := c.dg.NewTxn()
    ctx := context.Background()
    defer txn.Discard(ctx)
    queryStr := strings.Replace(data.QueryByUid, "$phone", uid, 1)
    logger.Infow("GetPersonByUid start", "uid", uid, "queryStr", queryStr)

    res, err := txn.Query(ctx, queryStr)
    if err != nil {
        logger.Panicw("GetPersonByUid execute query error", "uid", uid, "queryStr", queryStr, "err", err)
    }

    decode := DecodePerson{}
    err = json.Unmarshal(res.Json, &decode)
    if err != nil {
        logger.Panicw("GetPersonByUid unmarshal json  error", "uid", uid, "err", err, "res", string(res.Json))
    }

    if len(decode.P) == 0 {
        logger.Infow("GetPersonByUid node not exists", "uid", uid)
        return nil
    } else {
        logger.Infow("GetPersonByUid node exists", "uid", uid, "node", decode.P[0])
        return &decode.P[0]
    }
}

func (c *Client) GetPersonByEdge(edgeName string, value string) *data.Person {
    logger.Infow("GetPersonByEdge start", "edgeName", edgeName, "value", value)

    txn := c.dg.NewTxn()
    ctx := context.Background()
    defer txn.Discard(ctx)

    queryStr := strings.Replace(strings.Replace(data.QueryByEdge, "$edgeName", edgeName, 1), "$phone", value, 1)
    logger.Infow("GetPersonByEdge queryStr", "edgeName", edgeName, "value", value, "queryStr", queryStr)

    res, err := txn.Query(ctx, queryStr)
    if err != nil {
        logger.Panicw("GetPersonByEdge query error", "edgeName", edgeName, "value", value, "queryStr", queryStr, "err", err)
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

func (c *Client) DeleteByUid(uidList ...string) bool {
    logger.Infow("DeleteByUid start", "uidList", uidList)

    if len(uidList) == 0 {
        logger.Warnw("DeleteByUid uidList is empty", "uidList", uidList)
        return true
    }

    mutations := make([]string, len(uidList))
    for i, uid := range uidList {
        mutations[i] = strings.Replace(data.DeleteNode, "$uid", uid, 1)
    }
    deleteMutationList := strings.Join(mutations, "\n")
    logger.Warnw("DeleteByUid deleteMutationList", "deleteMutationList", deleteMutationList)

    mu := &api.Mutation{
        CommitNow: true,
        DelNquads: []byte(deleteMutationList),
    }
    txn := c.dg.NewTxn()
    ctx := context.Background()
    mutate, err := txn.Mutate(ctx, mu)
    if nil != err {
        logger.Panicw("DeleteByUid mutation error", "err", err)
    }

    logger.Info("DeleteByUid response", "response", string(mutate.Json))
    if mutate.Json == nil {
        return true
    }
    return false
}

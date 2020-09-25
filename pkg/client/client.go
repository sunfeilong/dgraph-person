package client

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
)

type Client struct {
	ip   string
	port int
	dg   *dgo.Dgraph
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

func getByUid(uid string, result *interface{}) {

}

func getByEdge(edgeName string, value string, result *interface{}) {

}

func deleteByUid(uid string) bool {
	return false
}

func deleteEdge(edgeName string, value string, result *interface{}) {

}

package data

type Person struct {
    Uid      string   `json:"uid,omitempty"`
    Name     string   `json:"name,omitempty"`
    Phone    string   `json:"phone,omitempty"`
    DType    []string `json:"dgraph.type,omitempty"`
    NodeType string   `json:"node.type,omitempty"`
}

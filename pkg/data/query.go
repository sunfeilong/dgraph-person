package data

const QueryByEdge string = `
  query {
    data(func: eq($edgeName, $phone)) {
      uid
      name
      phone
      node.type
      dgraph.type
    }
  }`

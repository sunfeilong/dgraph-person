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

const QueryByUid string = `
  query {
    data(func: uid($uid)) {
      uid
      name
      phone
      node.type
      dgraph.type
    }
  }`

const DeleteNode string = `<$uid> * * .`

const AddFriendRelation string = `<$uidA> <friend> <$uidB> .`

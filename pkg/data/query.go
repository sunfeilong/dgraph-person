package data

const QueryByEdge string = `
  query {
    data(func: eq($edgeName, $value)) {
      uid
      num
    }
  }`

const QueryByUid string = `
  query {
    data(func: uid($uid)) {
      uid
      num
    }
  }`

const DeleteNode string = `<$uid> * * .`

const AddFriendRelation string = `<$uidA> <friend> <$uidB> .`

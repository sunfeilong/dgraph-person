package data

const Schema string = `
  type Person{
    node.type # 用于表示节点的类型 (可能不需要)
    name      # 姓名
    phone     # 电话号码
    friend    # 朋友
  }

  name: string @index(term, exact) .
  phone: string  @index(exact, term) .
  friend: [uid] .
  node.type: string  @index(exact) .
  `

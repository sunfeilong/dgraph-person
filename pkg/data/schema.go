package data

const Schema string = `
  type Person{
    num       # 姓名
    friend    # 朋友
  }

  num: string @index(term, exact) .
  friend: [uid] .
  `

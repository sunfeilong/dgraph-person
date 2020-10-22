package reader

type PersonReader interface {
    ReadPersonFromFile(filePath string) []Num

    ReadFriendFromFile(filePath string) []NumAndFriendNum
}

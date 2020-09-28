package reader

type PersonReader interface {
    ReadPersonFromFile(filePath string) []IdNameAndPhone

    ReadFriendFromFile(filePath string) []PhoneNameAndPhone
}

package reader

type PersonReader interface {
    ReadFromFile(filePath string) *[]IdNameAndPhone
}

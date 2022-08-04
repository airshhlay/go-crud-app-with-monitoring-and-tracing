package db

type User struct {
	UserId   int64
	Username string
	Password []byte
}

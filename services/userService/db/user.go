package db

// User struct that defines the format of a user that is stored in a database.
type User struct {
	UserID   int64
	Username string
	Password []byte
}

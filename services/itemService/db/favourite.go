package db

import "time"

// Favourite is struct that defines the way a user's favourite is stored in the database. It follows the schema in schema/mysql.sql
type Favourite struct {
	ID        int64
	UserID    int64
	ItemID    int64
	ShopID    int64
	TimeAdded time.Time
}

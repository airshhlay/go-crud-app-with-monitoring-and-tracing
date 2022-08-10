package db

import "time"

type Favourite struct {
	Id        int64
	UserId    int64
	ItemId    int64
	ShopId    int64
	TimeAdded time.Time
}

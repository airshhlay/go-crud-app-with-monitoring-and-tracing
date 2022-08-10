package db

type Item struct {
	ItemId int    `redis:"itemId"`
	ShopId int    `redis:"shopId"`
	Price  int64  `redis:"price"`
	Name   string `redis:"name"`
}

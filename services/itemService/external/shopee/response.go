package shopee

type ShopeeGetItemRes struct {
	Error    int      `json:"error"`
	ErrorMsg string   `json:"error_msg"`
	ItemData ItemData `json:"data"`
}

type ItemData struct {
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	ItemId int64  `json:"itemId"`
	ShopId int64  `json:"shopId"`
}

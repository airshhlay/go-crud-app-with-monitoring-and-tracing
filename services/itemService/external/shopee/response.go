package shopee

// GetItemRes defines the expected response from the external HTTP call
// for the call to get an item's information
type GetItemRes struct {
	Error    int      `json:"error"`
	ErrorMsg string   `json:"error_msg"`
	ItemData ItemData `json:"data"`
}

// ItemData defines the required data in the GetItemRes response
type ItemData struct {
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	ItemID int64  `json:"itemid"`
	ShopID int64  `json:"shopid"`
}

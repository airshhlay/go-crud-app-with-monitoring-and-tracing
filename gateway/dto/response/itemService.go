package response

type AddFavRes struct {
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Item      Item   `json:"item"`
}

type Item struct {
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	ItemId int64  `json:"itemId"`
	ShopId int64  `json:"shopId"`
}

type DeleteFavRes struct {
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}

type GetFavListRes struct {
	ErrorCode  int32  `json:"errorCode"`
	ErrorMsg   string `json:"errorMsg"`
	Items      []Item `json:"items"`
	TotalPages int    `json:"totalPages"`
}

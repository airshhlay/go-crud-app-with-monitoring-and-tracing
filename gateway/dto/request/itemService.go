package request

type AddFavReq struct {
	ItemId int64 `json:"itemId"`
	ShopId int64 `json:"shopId"`
}

type DeleteFavReq struct {
	ItemId int64 `json:"itemId"`
	ShopId int64 `json:"shopId"`
}

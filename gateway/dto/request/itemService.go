package request

type AddFavReq struct {
	ItemId string `json:"itemId"`
	ShopId string `json:"shopId"`
}

type DeleteFavReq struct {
	ItemId int64 `json:"itemId"`
	ShopId int64 `json:"shopId"`
}

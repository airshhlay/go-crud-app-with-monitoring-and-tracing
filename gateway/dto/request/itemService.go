package request

// AddFavReq defines the expected request body to AddFav
type AddFavReq struct {
	ItemID string `json:"itemID"`
	ShopID string `json:"shopID"`
}

// DeleteFavReq defines the expected request body to DeleteFav
type DeleteFavReq struct {
	ItemID int64 `json:"itemID"`
	ShopID int64 `json:"shopID"`
}

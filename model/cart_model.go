package model

type InsertCartRequest struct {
	UserId    string
	Email     string
	ProductId string `json:"product_id"`
}

type CartResponse struct {
	Id           string  `json:"id"`
	ProductName  string  `json:"name"`
	ProductPrice float64 `json:"product_price"`
	ProductId    string  `json:"product_id"`
	Qty          int     `json:"qty"`
	TotalPrice   float64 `json:"total_price"`
}

type UpdateCartRequest struct {
	ProductId    string
	ProductName  string
	ProductPrice float64
	Qty          int
}

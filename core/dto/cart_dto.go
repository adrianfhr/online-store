package dto

type AddToCartDTO struct {
	ProductID  string `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required"`
}

type RemoveProductFromCartDTO struct {
	ProductID string `json:"product_id" binding:"required"`
}
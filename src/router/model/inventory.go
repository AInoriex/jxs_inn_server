package model

import (
	"time"
)

// @Title	获取用户藏品响应体
type GetInventoryListResponse struct {
	ProductId   string    `json:"product_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageUrl    string    `json:"image_url"`
	PurchaseAt  time.Time `json:"purchase_at"`
}

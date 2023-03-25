package model

type Transaction struct {
	ClientId  int64   `json:"client_id"`
	ServiceId int64   `json:"service_id"`
	OrderId   int64   `json:"order_id"`
	Price     float64 `json:"price"`
}

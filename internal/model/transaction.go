package model

type Transaction struct {
	ClientId  string `json:"client_id"`
	ServiceId string `json:"service_id"`
	OrderId   string `json:"order_id"`
	Price     string `json:"price"`
}

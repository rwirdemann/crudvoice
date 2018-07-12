package domain

type Project struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	CustomerId int    `json:"customerId"` // belongs to customer
}

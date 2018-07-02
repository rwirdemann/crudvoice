package domain

type Invoice struct {
	Id         int
	CustomerId int
	Month      int `json:"month"`
	Year       int `json:"year"`
	Status     string `json:"status"`
}

package domain

type Project struct {
	Id         int
	Name       string `json:"name"`
	CustomerId int
}

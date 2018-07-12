package domain

type Activity struct {
	Id     int
	Name   string `json:"name"`
	UserId string `json:"userId"` // belongs to user
}

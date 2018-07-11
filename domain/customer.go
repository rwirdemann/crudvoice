package domain

type Customer struct {
	Id     int
	Name   string
	UserId string // a user has many customers
}

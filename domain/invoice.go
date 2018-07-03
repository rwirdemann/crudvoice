package domain

import "io/ioutil"

type Invoice struct {
	Id         int
	CustomerId int
	Month      int    `json:"month"`
	Year       int    `json:"year"`
	Status     string `json:"status"`
	Bookings   []Booking
}

func (i *Invoice) Prepare() {
}

func (i Invoice) ToPdf() []byte {
	dat, _:= ioutil.ReadFile("/tmp/invoice.pdf")
	return dat
}

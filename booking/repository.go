package booking

import "github.com/rwirdemann/crudvoice/domain"

type Repository struct {
	bookings map[int]domain.Booking
}

func NewRepository() *Repository {
	return &Repository{bookings: make(map[int]domain.Booking)}
}

func (r *Repository) Create(booking domain.Booking) domain.Booking {
	booking.Id = r.nextId()
	r.bookings[booking.Id] = booking
	return booking
}

func (r *Repository) ByInvoiceId(invoiceId int) []domain.Booking {
	var bookings []domain.Booking
	for _, b := range r.bookings {
		if b.InvoiceId == invoiceId {
			bookings = append(bookings, b)
		}
	}
	return bookings
}

func (r *Repository) Delete(id int) {
	delete(r.bookings, id)
}

func (r *Repository) nextId() int {
	nextId := 1
	for _, v := range r.bookings {
		if v.Id >= nextId {
			nextId = v.Id + 1
		}
	}
	return nextId
}

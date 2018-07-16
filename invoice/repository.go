package invoice

import "github.com/rwirdemann/crudvoice/domain"

type Repository struct {
	invoices map[int]domain.Invoice
}

func NewRepository() *Repository {
	r := &Repository{invoices: make(map[int]domain.Invoice)}
	r.Create(domain.Invoice{Month: 6, Year: 2018, CustomerId: 1})
	return r
}

func (r *Repository) Create(invoice domain.Invoice) domain.Invoice {
	invoice.Id = r.nextId()
	invoice.Status = "open"
	invoice.Bookings = []domain.Booking{}
	r.invoices[invoice.Id] = invoice
	return invoice
}

func (r *Repository) Update(invoice domain.Invoice) {
	r.invoices[invoice.Id] = invoice
}

func (r *Repository) FindById(id int) (domain.Invoice, bool) {
	i, ok := r.invoices[id]
	return i, ok
}

func (r *Repository) nextId() int {
	nextId := 1
	for _, v := range r.invoices {
		if v.Id >= nextId {
			nextId = v.Id + 1
		}
	}
	return nextId
}

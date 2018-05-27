package database

import "github.com/rwirdemann/restvoice/domain"

type MySQLRepository struct {
	invoices []*domain.Invoice
}

func NewMySQLRepository() *MySQLRepository {
	r := MySQLRepository{}
	return &r
}

func (r *MySQLRepository) Invoices() []*domain.Invoice {
	return r.invoices
}

func (r *MySQLRepository) Create(invoice *domain.Invoice) {
	invoice.Id = r.nextInvoiceId()
	invoice.Status = "open"
	r.invoices = append(r.invoices, invoice)
}

func (r *MySQLRepository) nextInvoiceId() int {
	nextId := 1
	for _, i := range r.invoices {
			if i.Id >= nextId {
				nextId = i.Id + 1
		}
	}
	return nextId
}


package usecase

import "github.com/rwirdemann/restvoice/domain"

type Repository interface {
	Invoices() []*domain.Invoice
	Create(invoice *domain.Invoice)
	GetInvoice(id int) *domain.Invoice
}

package usecase

import (
	"github.com/rwirdemann/restvoice/foundation"
	"github.com/rwirdemann/restvoice/domain"
)

type CreateInvoice struct {
	invoiceConsumer foundation.Consumer
	repository      Repository
}

func NewCreateInvoice(invoiceConsumer foundation.Consumer, repository Repository) *CreateInvoice {
	return &CreateInvoice{
		repository:      repository,
		invoiceConsumer: invoiceConsumer}
}

func (u CreateInvoice) Run(i ...interface{}) interface{} {
	invoice := domain.NewInvoice()
	return invoice
}

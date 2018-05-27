package usecase

import (
	"github.com/rwirdemann/restvoice/foundation"
	"github.com/rwirdemann/restvoice/domain"
)

type CreateInvoice struct {
	consumer   foundation.Consumer
	presenter  foundation.Presenter
	repository Repository
}

func NewCreateInvoice(consumer foundation.Consumer, presenter foundation.Presenter, repository Repository) *CreateInvoice {
	return &CreateInvoice{
		repository: repository,
		consumer:   consumer,
		presenter:  presenter}
}

func (u CreateInvoice) Run(i ...interface{}) interface{} {
	invoice := u.consumer.Consume(i[0]).(*domain.Invoice)
	u.repository.Create(invoice)
	return u.presenter.Present(invoice)
}

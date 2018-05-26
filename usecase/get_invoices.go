package usecase

import (
	"github.com/rwirdemann/restvoice/foundation"
)

type GetInvoices struct {
	repository Repository
	presenter  foundation.Presenter
}

func NewGetInvoices(presenter foundation.Presenter, repository Repository) *GetInvoices {
	return &GetInvoices{presenter: presenter, repository: repository}
}

func (u GetInvoices) Run(i ...interface{}) interface{} {
	return u.presenter.Present(u.repository.Invoices())
}

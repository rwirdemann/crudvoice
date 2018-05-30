package usecase

import (
	"github.com/rwirdemann/restvoice/foundation"
	"strconv"
)

type GetInvoice struct {
	repository Repository
	consumer   foundation.Consumer
	presenter  foundation.Presenter
}

func NewGetInvoice(consumer foundation.Consumer, presenter foundation.Presenter, repository Repository) *GetInvoice {
	return &GetInvoice{consumer: consumer, presenter: presenter, repository: repository}
}

func (u GetInvoice) Run(i ...interface{}) interface{} {
	invoiceId, _ := strconv.Atoi(u.consumer.Consume(i[0]).(string))
	return u.presenter.Present(u.repository.GetInvoice(invoiceId))
}

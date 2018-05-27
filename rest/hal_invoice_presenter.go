package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rwirdemann/restvoice/domain"
)

type HALInvoice struct {
}

func NewHALInvoice() HALInvoice {
	return HALInvoice{}
}

func (j HALInvoice) present(i interface{}) interface{} {
	invoice := i.(*domain.Invoice)
	b, _ := json.Marshal(decorate(invoice))
	return b
}

type Link struct {
	Method string `json:"method"`
	Href   string `json:"href"`
}

type OperationsDecorator struct {
	*domain.Invoice
	Links map[string]Link `json:"_links"`
}

func decorate(i *domain.Invoice) OperationsDecorator {
	var links = make(map[string]Link)
	switch i.Status {
	case "open":
		links["self"] = Link{"GET", fmt.Sprintf("/invoice/%d", i.Id)}
		links["book"] = Link{"GET", fmt.Sprintf("/invoice/%d", i.Id)}
	}
	return OperationsDecorator{Invoice: i, Links: links}
}

func (j HALInvoice) Present(i interface{}) interface{} {
	var b []byte

	switch t := i.(type) {
	case []*domain.Invoice:
		var result []OperationsDecorator
		for _, i := range t {
			result = append(result, decorate(i))
		}
		b, _ = json.Marshal(result)
	case *domain.Invoice:
		b, _ = json.Marshal(decorate(t))
	}

	return b
}

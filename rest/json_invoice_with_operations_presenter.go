package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rwirdemann/restvoice/domain"
)

type JSONInvoiceWithOperations struct {
}

func NewJSONInvoiceWithOperationsPresenter() JSONInvoiceWithOperations {
	return JSONInvoiceWithOperations{}
}

func (j JSONInvoiceWithOperations) present(i interface{}) interface{} {
	invoice := i.(*domain.Invoice)
	b, _ := json.Marshal(decorate(invoice))
	return b
}

type Operation struct {
	Rel    string `json:"rel"`
	Method string `json:"method"`
	Href   string `json:"href"`
}

type OperationsDecorator struct {
	*domain.Invoice
	Operations []Operation `json:"operations"`
}

func decorate(i *domain.Invoice) OperationsDecorator {
	var operations []Operation
	switch i.Status {
	case "open":
		self := Operation{"self", "GET", fmt.Sprintf("/invoice/%d", i.Id)}
		book := Operation{"book", "POST", fmt.Sprintf("/booking/%d", i.Id)}
		operations = append(operations, self, book)
	}
	return OperationsDecorator{Invoice: i, Operations: operations}
}

func (j JSONInvoiceWithOperations) Present(i interface{}) interface{} {
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


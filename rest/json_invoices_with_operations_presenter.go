package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rwirdemann/restvoice/domain"
)

func NewJSONInvoicesWithOperationsPresenter() JSONInvoicesWithOperations {
	return JSONInvoicesWithOperations{}
}

type JSONInvoicesWithOperations struct {
}
type Operation struct {
	Rel    string `json:"rel"`
	Method string `json:"method"`
	Href   string `json:"href"`
}

type InvoiceWithOperations struct {
	domain.Invoice
	Operations []Operation `json:"operations"`
}

func (j JSONInvoicesWithOperations) Present(i interface{}) interface{} {
	var withOperations []InvoiceWithOperations
	invoices := i.([]domain.Invoice)
	for _, i := range invoices {
		show := Operation{"show", "GET", fmt.Sprintf("/invoices/%d", i.Id)}
		withOperation := InvoiceWithOperations{Invoice: i, Operations: []Operation{show}}
		withOperations = append(withOperations, withOperation)
	}

	b, _ := json.Marshal(withOperations)
	return b
}

package main

import (
	"github.com/rwirdemann/restvoice/database"
	"github.com/gorilla/mux"
	"fmt"
	"net/http"
	"github.com/rs/cors"
	"github.com/rwirdemann/restvoice/usecase"
	"github.com/rwirdemann/restvoice/domain"
	"github.com/rwirdemann/restvoice/rest"
)

func main() {
	invoiceConsumer := rest.NewJSONConsumer(&domain.Invoice{})
	pathVariableConsumer := rest.NewPathVariableConsumer("id")
	invoicesPresenter := rest.NewRVInvoicePresenter()
	invoicePresenter := rest.NewRVInvoicePresenter()

	repository := database.NewMySQLRepository()
	i := domain.NewInvoice("Libri GmbH")
	repository.Create(i)

	getInvoices := usecase.NewGetInvoices(invoicesPresenter, repository)
	getInvoice := usecase.NewGetInvoice(pathVariableConsumer, invoicesPresenter, repository)
	createInvoice := usecase.NewCreateInvoice(invoiceConsumer, invoicePresenter, repository)

	r := mux.NewRouter()
	r.HandleFunc("/invoice", rest.MakeGetInvoicesHandler(getInvoices)).Methods("GET")
	r.HandleFunc("/invoice/{id}", rest.MakeGetInvoiceHandler(getInvoice)).Methods("GET")
	r.HandleFunc("/invoice", rest.MakeCreateInvoiceHandler(createInvoice)).Methods("POST")

	fmt.Println("GET  http://localhost:8190/invoice")
	fmt.Println("GET  http://localhost:8190/invoice/1")
	fmt.Println("POST http://localhost:8190/invoice")

	http.ListenAndServe(":8190", cors.AllowAll().Handler(r))
}

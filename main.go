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
	invoicesPresenter := rest.NewHALInvoice()
	invoicePresenter := rest.NewHALInvoice()

	repository := database.NewMySQLRepository()
	i := domain.NewInvoice("Libri GmbH")
	repository.Create(i)

	getInvoices := usecase.NewGetInvoices(invoicesPresenter, repository)
	createInvoice := usecase.NewCreateInvoice(invoiceConsumer, invoicePresenter, repository)

	r := mux.NewRouter()
	r.HandleFunc("/invoice", rest.MakeGetInvoicesHandler(getInvoices)).Methods("GET")
	r.HandleFunc("/invoice", rest.MakeCreateInvoiceHandler(createInvoice)).Methods("POST")

	fmt.Println("GET  http://localhost:8190/invoice")
	fmt.Println("POST http://localhost:8190/invoice")

	http.ListenAndServe(":8190", cors.AllowAll().Handler(r))
}

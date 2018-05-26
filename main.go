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
	invoicePresenter := rest.NewJSONPresenter()
	repository := database.NewMySQLRepository()
	getInvoices := usecase.NewGetInvoices(invoicePresenter, repository)
	createInvoice := usecase.NewCreateInvoice(invoiceConsumer, repository)
	r := mux.NewRouter()
	r.HandleFunc("/invoices", rest.MakeGetInvoicesHandler(getInvoices)).Methods("GET")
	r.HandleFunc("/invoices", rest.MakeCreateInvoiceHandler(createInvoice)).Methods("POST")

	fmt.Println("GET  http://localhost:8190/invoices")
	fmt.Println("POST http://localhost:8190/invoices")

	http.ListenAndServe(":8190", cors.AllowAll().Handler(r))
}

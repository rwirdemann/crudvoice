package rest

import (
	"github.com/gorilla/mux"
	"net/http"
	"io/ioutil"
	"github.com/rwirdemann/restvoice/foundation"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	return r
}

func MakeCreateInvoiceHandler(usecase foundation.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		usecase.Run(body)
	}
}


package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/rwirdemann/crudvoice/domain"
	"github.com/rwirdemann/crudvoice/invoice"
	"strconv"
	"github.com/rwirdemann/crudvoice/booking"
	"time"
	"bytes"
	"os"
	"log"
	"github.com/joho/godotenv"
	"strings"
	"math/rand"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices",
		digestAuth(createInvoiceHandler)).Methods("POST")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}/bookings",
		basicAuth(createBookingHandler)).Methods("POST")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}/bookings/{bookingId:[0-9]+}", deleteBookingHandler).Methods("DELETE")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", updateInvoiceHandler).Methods("PUT")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", readInvoiceHandler).Methods("GET")

	http.ListenAndServe(":8080", r)
}

var invoiceRepository = invoice.NewRepository()
var bookingRepository = booking.NewRepository()

func createInvoiceHandler(writer http.ResponseWriter, request *http.Request) {
	// Read invoice data from request body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create invoice and marshal it to JSON
	var i domain.Invoice
	json.Unmarshal(body, &i)

	i.CustomerId, _ = strconv.Atoi(mux.Vars(request)["customerId"])
	created := invoiceRepository.Create(i)
	b, err := json.Marshal(created)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response
	writer.Header().Set("Location", fmt.Sprintf("%s/%d", request.URL.String(), created.Id))
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(b)
}

func createBookingHandler(writer http.ResponseWriter, request *http.Request) {
	// Read booking data from request body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create booking and marshal it to JSON
	var booking domain.Booking
	json.Unmarshal(body, &booking)
	created := bookingRepository.Create(booking)
	created.InvoiceId, _ = strconv.Atoi(mux.Vars(request)["invoiceId"])
	b, err := json.Marshal(created)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response
	writer.Header().Set("Location", fmt.Sprintf("%s/%d", request.URL.String(), created.Id))
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(b)
}

func deleteBookingHandler(writer http.ResponseWriter, request *http.Request) {
	bookingId, _ := strconv.Atoi(mux.Vars(request)["bookingId"])
	bookingRepository.Delete(bookingId)
	writer.WriteHeader(http.StatusNoContent)
}

func updateInvoiceHandler(writer http.ResponseWriter, request *http.Request) {
	// Read invoice data from request body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal and update invoice
	var i domain.Invoice
	json.Unmarshal(body, &i)
	i.Id, _ = strconv.Atoi(mux.Vars(request)["invoiceId"])
	i.CustomerId, _ = strconv.Atoi(mux.Vars(request)["customerId"])

	if i.Status == "payment expected" {
		i.Close()
	}
	invoiceRepository.Update(i)

	// Write response
	writer.WriteHeader(http.StatusNoContent)
}

func readInvoiceHandler(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(request)["invoiceId"])
	i, _ := invoiceRepository.FindById(id)
	accept := request.Header.Get("Accept")
	switch accept {
	case "application/pdf":
		http.ServeContent(writer, request, "invoice.pdf", time.Now(), bytes.NewReader(i.ToPdf()))
	case "application/json":
		b, _ := json.Marshal(i)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(b)
	default:
		writer.WriteHeader(http.StatusNotAcceptable)
	}
}

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if username, password, ok := r.BasicAuth(); ok {
			if username == os.Getenv("USERNAME") && password == os.Getenv("PASSWORD") {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", "Basic realm=\"example@restvoice.org\"")
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func digestAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := r.Header.Get("Authorization")
		if strings.HasPrefix(s, "Digest") {
			next.ServeHTTP(w, r)
			return
		}
		auth := fmt.Sprintf("Digest realm=\"example@restvoice.org\" qop=\"auth\" nonce=\"%s\" opaque=\"%s\"", nonce(), opaque())
		w.Header().Set("WWW-Authenticate", auth)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func nonce() string {
	return "UAZs1dp3wX5BtXEpoCXKO2lHhap564rX"
}

func opaque() string {
	return "xU2Z4FyqwKUBdwTMRYdGtAG1ppaT0bNm"
}

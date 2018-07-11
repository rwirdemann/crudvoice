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
	"crypto/md5"
	"io"
	"encoding/hex"
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"regexp"
	"github.com/rwirdemann/crudvoice/project"
	"github.com/rwirdemann/crudvoice/customer"
)

const validJWT = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.EkN-DOsnsuRjRO6BxXemmJDm3HbxrbRzXglbN2S4sOkopdU4IsDxTI8jO19W_A4K8ZPJijNLis4EZsHeY559a4DFOd50_OqgHGuERTqYZyuhtF39yxJPAjUESwxk2J5k_4zM3O-vtd1Ghyo4IbqKKSy6J9mTniYJPenn5-HIirE"

//const publicKeyFilePath = "valid_sample_key.pub"
const publicKeyFilePath = "keycloak_key.pub"

var publicKey *rsa.PublicKey

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if f, err := ioutil.ReadFile(publicKeyFilePath); err == nil {
		if publicKey, err = jwt.ParseRSAPublicKeyFromPEM(f); err != nil {
			log.Fatalf("Could not parse public key from pem file")
		}
	} else {
		log.Fatalf("Could not open public key file: %s", publicKeyFilePath)
	}

	r := mux.NewRouter()
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices",
		jwtAuth(assertCustomer(createInvoiceHandler))).Methods("POST")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}/bookings",
		digestAuth(createBookingHandler)).Methods("POST")

	r.HandleFunc("/customers/{customerId:[0-9]+}/projects",
		jwtAuth(assertAdmin(createProjectHandler))).Methods("POST")

	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}/bookings/{bookingId:[0-9]+}", deleteBookingHandler).Methods("DELETE")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", updateInvoiceHandler).Methods("PUT")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", readInvoiceHandler).Methods("GET")

	http.ListenAndServe(":8080", r)
}

var invoiceRepository = invoice.NewRepository()
var bookingRepository = booking.NewRepository()
var projectRepository = project.NewRepository()
var customerRepository = customer.NewRepository()

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

func jwtAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractJwtFromHeader(r.Header)
		if verifyJWT(token) {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", "Bearer realm=\"example@restvoice.org\"")
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func assertCustomer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractJwtFromHeader(r.Header)
		customerId, _ := strconv.Atoi(mux.Vars(r)["customerId"])
		if ownsCustomer(token, customerId) {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", "Bearer realm=\"example@restvoice.org\"")
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func ownsCustomer(token string, customerId int) bool {
	userId := claim(token, "sub")
	customer := customerRepository.ById(customerId)
	return customer.UserId == userId
}

func assertAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractJwtFromHeader(r.Header)
		if isAdmin(token) {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", "Bearer realm=\"example@restvoice.org\"")
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func claim(token string, key string) string {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err == nil {
		if claims, ok := t.Claims.(jwt.MapClaims); ok {
			if claims[key] != nil {
				return claims[key].(string)
			}
		}
	}

	return ""
}

func isAdmin(token string) bool {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err == nil {
		if claims, ok := t.Claims.(jwt.MapClaims); ok {
			if claims["admin"] != nil {
				return claims["admin"].(bool)
			}
		}
	}

	return false
}

func extractJwtFromHeader(header http.Header) (jwt string) {
	var jwtRegex = regexp.MustCompile(`^Bearer (\S+)$`)

	if val, ok := header["Authorization"]; ok {
		for _, value := range val {
			if result := jwtRegex.FindStringSubmatch(value); result != nil {
				jwt = result[1]
				return
			}
		}
	}

	return
}

var password = "time"

func digestAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if strings.HasPrefix(authorization, "Digest") {
			authFields := digestParts(authorization)
			step1 := hash(authFields["username"] + ":" + authFields["realm"] + ":" + password)
			step2 := hash(r.Method + ":" + authFields["uri"])
			step3 := hash(step1 + ":" +
				authFields["nonce"] + ":" +
				authFields["nc"] + ":" +
				authFields["cnonce"] + ":" +
				authFields["qop"] + ":" + step2)
			if step3 == authFields["response"] {
				next.ServeHTTP(w, r)
				return
			}
		}
		auth := fmt.Sprintf("Digest realm=\"%s\" qop=\"auth\" nonce=\"%s\" opaque=\"%s\"", realm(), nonce(), opaque())
		w.Header().Set("WWW-Authenticate", auth)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func hash(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

func digestParts(authorization string) map[string]string {
	result := map[string]string{}
	wantedHeaders := []string{"username", "nonce", "realm", "qop", "uri", "nc", "response", "opaque", "cnonce"}
	requestHeaders := strings.Split(authorization, ",")
	for _, r := range requestHeaders {
		for _, w := range wantedHeaders {
			if strings.Contains(r, " "+w) {
				v := strings.Split(r, "=")[1]
				result[w] = strings.Trim(v, `"`)
			}
		}
	}
	return result
}

func nonce() string {
	return "UAZs1dp3wX5BtXEpoCXKO2lHhap564rX"
}

func opaque() string {
	return "xU2Z4FyqwKUBdwTMRYdGtAG1ppaT0bNm"
}

func realm() string {
	return "example@restvoice.org"
}

func verifyJWT(token string) bool {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	return err == nil && t.Valid
}

func createProjectHandler(writer http.ResponseWriter, request *http.Request) {
	// Read invoice data from request body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create project and marshal it to JSON
	var p domain.Project
	json.Unmarshal(body, &p)

	p.CustomerId, _ = strconv.Atoi(mux.Vars(request)["customerId"])
	created := projectRepository.Create(p)
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

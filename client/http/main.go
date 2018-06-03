package main

import (
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/rwirdemann/restvoice/rest"
)

const baseUri = "http://localhost:8190/"

var client = &http.Client{}

func createInvoice(customer string, year int, month int) []byte {
	var jsonStr = []byte(fmt.Sprintf(`{"customer":"%s", "year": %d, "month": %d}`, customer, year, month))
	req, _ := http.NewRequest("POST", baseUri+"/invoice", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/vnd.restvoice+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

type Booking struct {
	Hours       float32 `json:"hours"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type Operations struct {
	Links map[string]rest.Link `json:"_links"`
}

var ops = Operations{Links: make(map[string]rest.Link)}

func book(hours float32, description string, date string) {
	b := Booking{Hours: hours, Description: description, Date: date}
	jsonStr, _ := json.Marshal(b)
	uri := baseUri + ops.Links["book"].Href
	req, _ := http.NewRequest(ops.Links["book"].Method, uri, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/vnd.restvoice+json")
	resp, _ := client.Do(req);
	defer resp.Body.Close()
	fmt.Printf("Status: %s", resp.Status)
}

func main() {
	body := createInvoice("Volkswagen AG", 2018, 8)
	json.Unmarshal(body, &ops)
	book(4.5, "NRG-333 Benutzersynchronisation", "2018-01-17")
}

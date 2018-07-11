package main

import (
	"net/http"
	"crypto/x509"
	"crypto/tls"
	"io/ioutil"
	"fmt"
)

func main() {
	certPool := x509.NewCertPool()
	rootCert, _ := ioutil.ReadFile("localhost.crt")
	certPool.AppendCertsFromPEM(rootCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: certPool},
		},
	}

	r, _ := client.Get("https://localhost:8081/hello")
	fmt.Printf("Status: %s", r.Status)
}
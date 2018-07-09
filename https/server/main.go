package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/hello", hello)
	http.ListenAndServeTLS(":8081", "localhost.crt", "localhost.key", nil)
}

func hello(writer http.ResponseWriter, _ *http.Request) {
	writer.Write([]byte("hello"))
}

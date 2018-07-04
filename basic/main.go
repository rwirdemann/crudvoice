package main

import (
	b64 "encoding/base64"
	"fmt"
)

func main() {
	print(b64.StdEncoding.EncodeToString([]byte("go:time")))
	sDec, _ := b64.StdEncoding.DecodeString("Z286dGltZQ==")
	fmt.Println(string(sDec))
}

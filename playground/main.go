package main

import "fmt"

func main() {
	m := make(map[string]interface{})

	m["user"] = true

	fmt.Printf("Value: %v", m["admin"].(bool))
}

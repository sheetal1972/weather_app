package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server listening on :8080...")
	http.ListenAndServe(":8080", nil)
}

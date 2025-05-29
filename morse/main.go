package main

import (
	"log"
	"net/http"

	handler "morse/api"
)

func main() {
	fs := http.FileServer(http.Dir("./public/"))
	http.Handle("/", fs)

	http.HandleFunc("/api/convert", handler.Handler)
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

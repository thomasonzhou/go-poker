package main

import (
	"log"
	"net/http"
)

func main() {

	store := NewInMemoryStore()
	server := NewPlayerServer(store)
	log.Fatal(http.ListenAndServe(":8000", server))
}

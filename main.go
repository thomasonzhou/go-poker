package main

import (
	"log"
	"net/http"
)

func main() {
	store := NewInMemoryStore()
	server := &PlayerServer{store}
	log.Fatal(http.ListenAndServe(":8000", server))
}

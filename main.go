package main

import (
	"fmt"
	"go-PostgresSql/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()

	fmt.Printf("Starting server on the port 8080....")

	log.Fatal(http.ListenAndServe(":8080", r))
}

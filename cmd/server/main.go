package main

import (
	"log"

	"github.com/dmehra2102/go-distributed-system/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
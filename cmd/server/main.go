package main

import (
	"log"

	"github.com/sprectza/proglog/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatalf(srv.ListenAndServe().Error())
}

package main

import (
	"github.com/jordanharrington/playground/goblockgo/internal/handler"
	"github.com/jordanharrington/playground/goblockgo/internal/service"
	"log"
	"net/http"
)

func main() {
	s := service.NewSimpleBlockchainService()

	log.Fatal(http.ListenAndServe(":8080", handler.Route(s)))
}

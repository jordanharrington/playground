package main

import (
	"github.com/jordanharrington/playground/goblockgo/internal/handler"
	"github.com/jordanharrington/playground/goblockgo/internal/repl"
	"github.com/jordanharrington/playground/goblockgo/internal/service"
	"log"
	"net/http"
	"os"
)

func main() {
	s := service.NewSimpleBlockchainService()

	mode := os.Getenv("MODE")
	switch mode {
	case "INTERACTIVE":
		finish := make(chan bool)
		// Serve http request on port 8080
		go listen(s, ":8080")
		// Start REPL
		go startRepl(s, finish)
		// Wait for REPL to exit
		<-finish
	default:
		listen(s, ":8080")
	}
}

func listen(service service.GoBlockGo, port string) {
	log.Fatal(http.ListenAndServe(port, handler.Route(service)))
}

func startRepl(service service.GoBlockGo, finish chan<- bool) {
	repl.Eval(service)
	finish <- true
}

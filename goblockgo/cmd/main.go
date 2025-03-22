package main

import (
	"fmt"
	"github.com/jordanharrington/playground/goblockgo/internal/handler"
	"github.com/jordanharrington/playground/goblockgo/internal/repl"
	"github.com/jordanharrington/playground/goblockgo/internal/service"
	"github.com/mattn/go-isatty"
	"log"
	"net/http"
	"os"
)

func main() {
	s := service.NewSimpleBlockchainService()
	// Serve http request on port 8080
	go listen(s, ":8080")
	// Start REPL
	finish := make(chan bool)
	go startRepl(s, finish)
	// Wait for REPL to exit
	<-finish
}

func listen(service service.GoBlockGo, port string) {
	log.Fatal(http.ListenAndServe(port, handler.Route(service)))
}

func startRepl(service service.GoBlockGo, finish chan<- bool) {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Println("stdin is not a terminal — skipping REPL")
		return
	}

	repl.Eval(service)
	finish <- true
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	gobootweb "github.com/furkilic/go-boot-web/pkg/go-boot-web"
)

func main() {
	errCh := gobootweb.Start()

	log.Println("Started")
	defer gobootweb.Stop()

	router := gobootweb.Router()
	router.Methods("GET").Path("/hello").Name("Hello").HandlerFunc(hello)

	log.Println("Routed")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Waiting")
	select {
	case err := <-errCh:
		log.Println(err)
		break
	case err := <-termChan:
		log.Println(err)
		break
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Received Request")
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Hello World !!")
}

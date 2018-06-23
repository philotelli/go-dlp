package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags

	r := mux.NewRouter()

	r.Handle("/lookup", &lookupHandler{}).Methods("GET", "POST", "DELETE")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(*addr, r))
}

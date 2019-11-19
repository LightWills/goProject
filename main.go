package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", helloWord).Methods("GET")
	r.HandleFunc("/bye", goodbye).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func helloWord(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Choupina")
}

func goodbye(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "By Choupina ")
}

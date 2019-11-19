package main

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var client *redis.Client
var templates *template.Template

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	templates = template.Must(templates.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	val, err := client.LRange("comments", 0, 10).Result()
	if err != nil {
		panic(err)
	}
	templates.ExecuteTemplate(w, "index.html", val)

}

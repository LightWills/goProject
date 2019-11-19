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
	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/", indexPostHandler).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	val, err := client.LRange("comments", 0, 10).Result()
	if err != nil {
		panic(err)
	}
	templates.ExecuteTemplate(w, "index.html", val)

}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment_element := r.PostForm.Get("comment_element")
	client.LPush("comments", comment_element)
	http.Redirect(w, r, "/", 302)

}

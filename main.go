package main

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
)

var client *redis.Client
var templates *template.Template
var store = sessions.NewCookieStore([]byte("kjsdjc7675yjhbks"))

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	templates = template.Must(templates.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/", indexPostHandler).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/test", testGetHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
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

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)

}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("name")
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)

	templates.ExecuteTemplate(w, "login.html", nil)

}

func testGetHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	untyped, ok := session.Values["username"]
	if !ok {
		return
	}
	username, ok := untyped.(string)
	if !ok {
		return
	}
	w.Write([]byte(username))
}

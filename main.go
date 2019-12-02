package main

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
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
	r.HandleFunc("/", AuthRequired(indexGetHandler)).Methods("GET")
	r.HandleFunc("/", AuthRequired(indexPostHandler)).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	r.HandleFunc("/test", testGetHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
		}

		handler.ServeHTTP(w, r)
	}
}
func indexGetHandler(w http.ResponseWriter, r *http.Request) {

	val, err := client.LRange("comments", 0, 10).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error - Index"))
		return
	}
	templates.ExecuteTemplate(w, "index.html", val)

}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment_element := r.PostForm.Get("comment_element")
	err := client.LPush("comments", comment_element)
	if err == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(w, r, "/", 302)

}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)

}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("name")
	pwd := r.PostForm.Get("password")
	session, _ := store.Get(r, "session")
	var k = client.Get("user: " + username)
	hash, err := k.Bytes()
	if err == redis.Nil {
		templates.ExecuteTemplate(w, "login.html", "unknow user")
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(pwd))
	if err != nil {
		templates.ExecuteTemplate(w, "login.html", "invalid login")
		return
	}
	else {
		session, _ = store.Get(r, "session")
		session.Values["username"] = username
		session.Save(r, w)
		http.Redirect(w, r, "/", 302)

		return
	}



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
func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", nil)

}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	firstname := r.PostForm.Get("firstname")
	password := r.PostForm.Get("password")
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	err = client.Set("user: "+firstname, hash, 0).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	/*session, _ := store.Get(r, "session")
	session.Values["username"] = firstname
	session.Save(r, w)*/

	//templates.ExecuteTemplate(w, "login.html", nil)
	http.Redirect(w, r, "/login", 302)

}

package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var urlsMap map[string]string
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")

// Page contains the data for template rendering
type Page struct {
	// Url is the shortened URL
	URL string
}

func randString(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	urls, ok := query["url"]

	if !ok || len(urls[0]) < 1 {
		log.Println("Nothing to do.")
		return
	}

	randStr := randString(32)
	urlsMap[randStr] = urls[0]

	p := Page{URL: randStr}

	t, err := template.ParseFiles("result.html")
	if err != nil {
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Cannot load the template properly"))
			log.Printf("Cannot open the template file, errored: %s", err.Error())
		}
	}

	t.Execute(w, p)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	log.Printf("Processing query: %s", r.URL)
	qs, ok := r.URL.Query()["q"]

	if !ok || len(qs[0]) < 1 {
		log.Println("Nothing to do.")
		return
	}

	log.Printf("Redirecting to: %s", urlsMap[qs[0]])
	w.WriteHeader(303)
	http.Redirect(w, r, urlsMap[qs[0]], http.StatusSeeOther)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tiny.html")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Cannot load the template properly"))
		log.Printf("Cannot open the template file, errored: %s", err.Error())
	}

	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/r", handleRedirect)
	http.HandleFunc("/", handleMain)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init() {
	rand.Seed(time.Now().UnixNano())
	urlsMap = make(map[string]string)
}

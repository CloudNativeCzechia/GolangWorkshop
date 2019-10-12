package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Storage represent url shortener storage
type Storage interface {
	Set(short, full string)
	Get(short string) (string, error)
}

type urlStorage struct {
	urlsMap map[string]string
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")

// Page contains the data for template rendering
type Page struct {
	// Url is the shortened URL
	URL string
}

var storage urlStorage

func (s *urlStorage) Set(short, full string) {
	s.urlsMap[short] = full
}

func (s *urlStorage) Get(short string) (string, error) {
	full, ok := s.urlsMap[short]
	if !ok {
		return "", fmt.Errorf("No key found")
	}
	return full, nil
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
	storage.Set(randStr, urls[0])

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

	fullURL, err := storage.Get(qs[0])
	if err != nil {
		log.Print("No url found in store")
	}
	log.Printf("Redirecting to: %s", fullURL)
	w.WriteHeader(303)
	http.Redirect(w, r, fullURL, http.StatusSeeOther)
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
	storage.urlsMap = make(map[string]string)
}

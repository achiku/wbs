package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	log.Print(url.QueryEscape("Hello, World"))
	fmt.Fprint(w, "Hello, world!\n")
}

func main() {
	helloWorldHandler := http.HandlerFunc(helloWorld)
	http.Handle("/hello", loggingMiddleware(helloWorldHandler))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

package main

// build upon: https://medium.com/@wembleyleach/simple-web-application-with-go-24ba8acf4c1

import (
	"bookfund/transaction"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func IsDevelop() bool {
	return os.Getenv("APP_ENV") == "develop"
}

// logging is middleware for wrapping any handler we want to track response
// times for and to see what resources are requested.
func LoggingHandler(handler func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		req := fmt.Sprintf("%s %s", r.Method, r.URL)
		log.Println(req)
		handler(w, r)
		log.Println(req, "completed in", time.Since(start))
	})
}

func main() {
	log.Print("Starting...")

	log.Print("Setting up environment.")

	log.Print("Connecting to database.")
	dbName := "bookfund"
	if IsDevelop() {
		dbName += "-develop"
	}
	err := transaction.OpenDB(fmt.Sprintf("./db/%s.db", dbName))
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer transaction.CloseDB()

	mux := http.NewServeMux()
	// public serves static assets such as CSS and JavaScript to clients.
	mux.Handle("GET /public/", LoggingHandler(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))).ServeHTTP(w, r)
	}))
	mux.Handle("GET /robots.txt", http.FileServer(http.Dir(".")))
	mux.Handle("GET /", LoggingHandler(index))
	mux.Handle("POST /{type}", LoggingHandler(post))
	mux.Handle("DELETE /{id}", LoggingHandler(delete))

	mux.Handle("GET /modal/{type}", LoggingHandler(modal))

	mux.Handle("GET /review", LoggingHandler(review))
	mux.Handle("GET /review/search", LoggingHandler(reviewSearch))

	mux.Handle("GET /login", LoggingHandler(login))
	mux.Handle("GET /logout", LoggingHandler(logout))

	port := os.Getenv("PORT")
	addr := fmt.Sprintf("localhost:%s", port)
	server := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("Done.")
	log.Printf("Listening on port %s (http://%s)\n", port, addr)
	log.Fatal(server.ListenAndServe())
}

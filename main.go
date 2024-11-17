package main

// build upon: https://medium.com/@wembleyleach/simple-web-application-with-go-24ba8acf4c1

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"database/sql"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

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

var db *sql.DB

func main() {
	log.Print("Starting...")

	log.Print("Setting up environment.")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Print("Connecting to database.")
	dbName := fmt.Sprintf("./db/%s.db", os.Getenv("DB_NAME"))
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	// public serves static assets such as CSS and JavaScript to clients.
	mux.Handle("GET /public/", LoggingHandler(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))).ServeHTTP(w, r)
	}))
	mux.Handle("GET /", LoggingHandler(index))
	mux.Handle("GET /review", LoggingHandler(review))
	mux.Handle("GET /modal/{type}", LoggingHandler(modal))
	mux.Handle("POST /transaction/{type}", LoggingHandler(transaction))

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

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/takshpanchal/crud-go/handlers"
	"github.com/takshpanchal/crud-go/models"
)

func main() {
	// Loading Env variable
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setup Loggers
	infoLogger := log.New(os.Stdin, "INFO: ", log.Ldate|log.Ltime)
	errLogger := log.New(os.Stdin, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)

	port := os.Getenv("PORT")

	// DB Setup
	db, err := setupPsqlDB()
	if err != nil {
		errLogger.Fatal(err)
	}

	//handlers
	mux := http.NewServeMux()
	userModel := &models.UserModel{DB: db}
	userHandler := &handlers.UserHandler{InfoLogger: infoLogger, ErrLogger: errLogger, Model: userModel}

	mux.HandleFunc("/user/create", userHandler.CreateUser)
	mux.HandleFunc("/user/", userHandler.User)

	infoLogger.Printf("Starting server at port %s", port)
	srv := &http.Server{Addr: port, Handler: mux, ErrorLog: errLogger}
	errLogger.Fatal(srv.ListenAndServe())
}

func setupPsqlDB() (*sql.DB, error) {
	host := os.Getenv("PGHOST")
	user := os.Getenv("PGUSER")
	dbname := os.Getenv("PGDATABASE")
	pass := os.Getenv("PGPASSWORD")
	connStr := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=verify-full", host, user, dbname, pass)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// Pinging DB to check for connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

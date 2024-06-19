package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"weddi.org/vendor-api/internal/database"
)

type config struct {
	port int
	env string
}

type apiConfig struct {
	config config
	DB *database.Queries
	apiKey string
	logger *log.Logger
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "DEVELOPMENT", "Environment (DEVELOPMENT|STAGING|PRODUCTION)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)

	_ = godotenv.Load()
	apiKey := os.Getenv("API_KEY")
	dbUrl := os.Getenv("DB_CONN")
	print(dbUrl)
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		config: cfg,
		DB: dbQueries,
		apiKey: apiKey,
		logger: logger,
	}
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/v1").Subrouter()

	subrouter.HandleFunc("/reviews", apiCfg.handleGetVendorReviews).Methods("GET")
	subrouter.HandleFunc("/reviews", apiCfg.handleCreateVendorReview).Methods("POST")
	subrouter.HandleFunc("/reviews/vendors/{id}", apiCfg.handleGetVendorReviewByVendorId).Methods("GET")
	subrouter.HandleFunc("/reviews/{id}", apiCfg.handleDeleteVendorReview).Methods("DELETE")
	subrouter.HandleFunc("/reviews/{id}", apiCfg.handleUpdateVendorReview).Methods("PATCH")	
	subrouter.HandleFunc("/scrape-reviews/{id}", apiCfg.handleScrapeReviewImage).Methods("POST")
	subrouter.HandleFunc("/list-images", apiCfg.handleListImagesByVendorId).Methods("GET")
	subrouter.HandleFunc("/review-images", apiCfg.handleGetImageByReviewId).Methods("GET")

	server := &http.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", cfg.port),
		Handler: router,
		IdleTimeout: time.Minute,
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on %s", cfg.env, server.Addr)
	err = server.ListenAndServe()
	logger.Fatal(err)
}

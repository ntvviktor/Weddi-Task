package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"weddi.org/vendor-api/config"
	controller "weddi.org/vendor-api/controller/vendor_review"
	"weddi.org/vendor-api/middleware"
)

func main() {
	var apiCfg config.ApiConfig
	flag.IntVar(&apiCfg.Port, "port", 8080, "API server port")
	flag.StringVar(&apiCfg.Env, "env", "DEVELOPMENT", "Environment (DEVELOPMENT|STAGING|PRODUCTION)")
	flag.Parse()

	_ = godotenv.Load()
	apiKey := os.Getenv("API_KEY")
	apiCfg.ApiKey = apiKey

	config.DB = config.SetupDatabase()
	
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/v1").Subrouter()

	subrouter.Use(middleware.LoggingMiddleware)
	subrouter.HandleFunc("/reviews", controller.GetVendorReviewsHandler).Methods("GET")
	subrouter.HandleFunc("/reviews", controller.CreateVendorReviewHandler).Methods("POST")
	subrouter.HandleFunc("/reviews/vendors/{id}", controller.GetVendorReviewByVendorIdHandler).Methods("GET")
	subrouter.HandleFunc("/reviews/{id}", controller.DeleteVendorReviewHandler).Methods("DELETE")
	subrouter.HandleFunc("/reviews/{id}", controller.UpdateVendorReviewHandler).Methods("PATCH")	
	subrouter.HandleFunc("/scrape-reviews/{id}", controller.ScrapeReviewImageHandler).Methods("POST")
	subrouter.HandleFunc("/vendor-images", controller.GetImagesByVendorIdHandler).Methods("GET")
	subrouter.HandleFunc("/review-images",controller.GetImageByReviewIdHandler).Methods("GET")

	server := &http.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", apiCfg.Port),
		Handler: router,
		IdleTimeout: time.Minute,
		ReadTimeout: 30 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	fmt.Printf("Starting %s server on %s \n", apiCfg.Env, server.Addr)
	err := server.ListenAndServe()
	log.Fatal(err)
}

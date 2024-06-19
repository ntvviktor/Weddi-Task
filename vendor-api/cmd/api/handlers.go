package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"weddi.org/vendor-api/internal/database"
)

type VendorReview struct {
	ReviewID 	 string `json:"review_id,omitempty"`
	VendorID     string `json:"vendor_id"`
	VendorName   string `json:"vendor_name"`
	Poster       string `json:"poster"`
	Date         string `json:"date"`
	Rating       int32  `json:"rating"`
	Source       string `json:"source"`
	Content      string `json:"content"`
	LinkToSource string `json:"link_to_source"`
}

func (apiConfig *apiConfig) handleGetVendorReviews(w http.ResponseWriter, req *http.Request) {
	data, err := apiConfig.DB.GetAllVendorReviews(req.Context())
	if err != nil {
		apiConfig.logger.Fatal(err)
		http.Error(w, "Interal error from database", http.StatusInternalServerError)
		return
	}
	var convertedData []VendorReview
	for _, vendorReview := range data {
		convertedData = append(convertedData, VendorReview{
            ReviewID:     vendorReview.ReviewID,
            VendorID:     vendorReview.VendorID,
            VendorName:   vendorReview.VendorName,
            Poster:       vendorReview.Poster,
            Date:         vendorReview.Date,
            Rating:       vendorReview.Rating,
            Source:       vendorReview.Source,
            Content:      vendorReview.Content,
            LinkToSource: vendorReview.LinkToSource,
        })
	}

	vendorReviews, err := json.Marshal(convertedData)
	if err != nil {
		apiConfig.logger.Fatal(err)
		http.Error(w, "Interal error in parsing json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(vendorReviews)
}

func (apiConfig *apiConfig) handleGetVendorReviewByVendorId(w http.ResponseWriter, req *http.Request) {
	vendorId := mux.Vars(req)["id"]
	vendorReviews, err := apiConfig.getVendorReviewByVendorId(req, vendorId)
	if err != nil {
		http.Error(w, "Interal error from database", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(vendorReviews)
}

func (apiConfig *apiConfig) handleCreateVendorReview(w http.ResponseWriter, req *http.Request) {
	var params VendorReview
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Bad json request", http.StatusBadRequest)
		return
	}

	reviewId := uuid.New()

	newVendorReview := database.CreateVendorReviewParams{
		ReviewID     : reviewId.String(),
		VendorID     : params.VendorID,
		VendorName   : params.VendorName,
		Poster       : params.Poster,
		Date         : params.Date,
		Rating       : params.Rating,
		Source       : params.Source,
		Content      : params.Content,
		LinkToSource : params.LinkToSource,
	}
	
	err = apiConfig.DB.CreateVendorReview(req.Context(), newVendorReview)
	if err != nil {
		apiConfig.logger.Fatal(err)
		http.Error(w, "Interal error from database", http.StatusInternalServerError)
		return
	}

	VendorReview, err := apiConfig.getVendorReviewByReviewId(req, reviewId.String())
	if err != nil {
		apiConfig.logger.Fatal(err)
		http.Error(w, "Interal error from database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(VendorReview)
}

func (apiConfig *apiConfig) handleUpdateVendorReview(w http.ResponseWriter, req *http.Request) {
	var params struct {
		VendorName   string `json:"vendor_name"`
		Poster       string `json:"poster"`
		Date         string `json:"date"`
		Rating       int32  `json:"rating"`
		Source       string `json:"source"`
		Content      string `json:"content"`
		LinkToSource string `json:"link_to_source"`
	}
	reviewId := mux.Vars(req)["id"]
	_, err := apiConfig.getVendorReviewByReviewId(req, reviewId)
	if err != nil {
		http.Error(w, "No resource match your review id", http.StatusNotFound)
		return
	}
	err = json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Bad json request", http.StatusBadRequest)
		return
	}

	newVendorReview := database.UpdateVendorReviewParams{
		VendorName   : params.VendorName,
		Poster       : params.Poster,
		Date         : params.Date,
		Rating       : params.Rating,
		Source       : params.Source,
		Content      : params.Content,
		LinkToSource : params.LinkToSource,
		ReviewID: reviewId,
	}

	b, err := json.MarshalIndent(newVendorReview, "", "  ")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Print(string(b))
	
	err = apiConfig.DB.UpdateVendorReview(req.Context(), newVendorReview)
	if err != nil {
		apiConfig.logger.Fatal(err)
		http.Error(w, "Unable to update! Interal error from database", http.StatusInternalServerError)
		return
	}

	VendorReview, err := apiConfig.getVendorReviewByReviewId(req, reviewId)
	if err != nil {
		apiConfig.logger.Fatal(err)
		http.Error(w, "Interal error from database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(VendorReview)
}

func (apiConfig *apiConfig) handleDeleteVendorReview(w http.ResponseWriter, req *http.Request) {
	vendorId := mux.Vars(req)["id"]
	err := apiConfig.DB.DeleteVendorReviewByReviewId(req.Context(), vendorId)
	if err != nil {
		apiConfig.logger.Fatal(err)
		http.Error(w, "Interal error from database", http.StatusInternalServerError)
		return	
	}
	w.WriteHeader(http.StatusOK)
}

func (apiConfig *apiConfig) handleScrapeReviewImage(w http.ResponseWriter, req *http.Request) {
	_ = godotenv.Load()
	scrapeUrl := os.Getenv("SCRAPE_URL")
	reviewId := mux.Vars(req)["id"]
	data, err := apiConfig.DB.GetVendorReviewByReviewId(req.Context(), reviewId)
	if err != nil {
		http.Error(w, "No resource match your review id", http.StatusNotFound)
		return
	}

	reviewUrl := data.LinkToSource

	resp, err := scrapeReviewVendorImage(scrapeUrl, reviewId, reviewUrl)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return	
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}


func (apiConfig *apiConfig) handleListImagesByVendorId(w http.ResponseWriter, req *http.Request) {
	type ImageResponse struct {
		ReviewID string `json:"review_id"`
		Images []string `json:"images"`
	}
	vendorId := req.URL.Query().Get("id")
    if vendorId == "" {
        http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
        return
    }
	
	// Get the vendor id from the database
	vendorReviews, err := apiConfig.DB.GetVendorReviewByVendorId(req.Context(), vendorId)
	if err != nil {
		http.Error(w, "Vendor ID not found", http.StatusNotFound)
        return
	}

	var reviewIds []string
	for _, review := range vendorReviews {
		reviewIds = append(reviewIds, review.ReviewID)
	}
	var response []ImageResponse
	for _, reviewId := range reviewIds {

		// Construct the directory path
		imagesDir := filepath.Join("..", "images", reviewId)
	
		// Read the directory
		files, err := os.ReadDir(imagesDir)
		if err != nil {
			fmt.Println("Image path not exist")
			continue
		}
	
		var imageFiles []string
		for _, file := range files {
			if !file.IsDir() {
				imageFiles = append(imageFiles, filepath.Join("/images", reviewId, file.Name()))
			}
		}
		imageRes := ImageResponse{
			ReviewID: reviewId,
			Images: imageFiles,
		}

		response = append(response, imageRes)
	}
    
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Error generating JSON response: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}


func (apiConfig *apiConfig) handleGetImageByReviewId(w http.ResponseWriter, req *http.Request) {
	reviewId := req.URL.Query().Get("review-id")
	imgId := req.URL.Query().Get("image-id")

	var Path = fmt.Sprintf("../images/%s/%s.png", reviewId, imgId) 

    img, err := os.Open(Path)
    if err != nil {
        apiConfig.logger.Fatal(err)
		http.Error(w, "Image not found", http.StatusNotFound)
		return
    }
    defer img.Close()
    w.Header().Set("Content-Type", "image/png")
    io.Copy(w, img)
}
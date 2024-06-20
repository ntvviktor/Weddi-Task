package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"weddi.org/vendor-api/config"
	"weddi.org/vendor-api/internal/database"
	"weddi.org/vendor-api/utils"
)

func GetVendorReviewsHandler(w http.ResponseWriter, req *http.Request) {
	data, err := config.DB.GetAllVendorReviews(req.Context())
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Interal error from database")
		log.Fatal(err)
		return
	}
	var convertedData []utils.VendorReviewHelper
	for _, vendorReview := range data {
		convertedData = append(convertedData, utils.VendorReviewHelper{
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
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Interal error when parsing json")
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(vendorReviews)
}

func GetVendorReviewByVendorIdHandler(w http.ResponseWriter, req *http.Request) {
	vendorId := mux.Vars(req)["id"]
	vendorReviews, err := utils.GetVendorReviewByVendorId(req, vendorId)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Interal error from database")
		log.Fatal(err)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(vendorReviews)
}

func CreateVendorReviewHandler(w http.ResponseWriter, req *http.Request) {
	var params utils.VendorReviewHelper
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Bad JSON request")
		log.Fatal(err)
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
	
	err = config.DB.CreateVendorReview(req.Context(), newVendorReview)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Interal error from database")
		log.Fatal(err)
		return
	}

	VendorReview, err := utils.GetVendorReviewByReviewId(req, reviewId.String())
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Interal error from database")
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(VendorReview)
}

func  UpdateVendorReviewHandler(w http.ResponseWriter, req *http.Request) {
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
	_, err := utils.GetVendorReviewByReviewId(req, reviewId)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusNotFound, "No resource match your review id")
		log.Fatal(err)
		return
	}
	err = json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusNotFound, "Bad JSON request")
		log.Fatal(err)
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

	// b, err := json.MarshalIndent(newVendorReview, "", "  ")
    // if err != nil {
    //     fmt.Println(err)
    // }
    // fmt.Print(string(b))
	
	err = config.DB.UpdateVendorReview(req.Context(), newVendorReview)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Unable to update! Interal error from database")
		log.Fatal(err)
		return
	}

	VendorReview, err := utils.GetVendorReviewByReviewId(req, reviewId)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Interal error from database")
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(VendorReview)
}

func DeleteVendorReviewHandler(w http.ResponseWriter, req *http.Request) {
	vendorId := mux.Vars(req)["id"]
	err := config.DB.DeleteVendorReviewByReviewId(req.Context(), vendorId)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Interal error from database")
		log.Fatal(err)
		return	
	}
	w.WriteHeader(http.StatusOK)
}

func ScrapeReviewImageHandler(w http.ResponseWriter, req *http.Request) {
	_ = godotenv.Load()
	scrapeUrl := os.Getenv("SCRAPE_URL")
	reviewId := mux.Vars(req)["id"]
	data, err := config.DB.GetVendorReviewByReviewId(req.Context(), reviewId)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusNotFound, "No resource match your review id")
		return
	}

	reviewUrl := data.LinkToSource

	resp, err := utils.ScrapeReviewVendorImage(scrapeUrl, reviewId, reviewUrl)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Interal error from database")
		return	
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}


func GetImagesByVendorIdHandler(w http.ResponseWriter, req *http.Request) {
	type ImageResponse struct {
		ReviewID string `json:"review_id"`
		Images []string `json:"images"`
	}
	vendorId := req.URL.Query().Get("id")
    if vendorId == "" {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Missing 'id' query parameter")
        return
    }
	
	// Get the vendor id from the database
	vendorReviews, err := config.DB.GetVendorReviewByVendorId(req.Context(), vendorId)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusNotFound, "Vendor ID not found")
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
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Error generating JSON response: ")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}


func GetImageByReviewIdHandler(w http.ResponseWriter, req *http.Request) {
	reviewId := req.URL.Query().Get("review-id")
	imgId := req.URL.Query().Get("image-id")

	var Path = fmt.Sprintf("../images/%s/%s.png", reviewId, imgId) 

    img, err := os.Open(Path)
    if err != nil {
		utils.NewErrorResponse(w, http.StatusNotFound, "Image not found")
        log.Fatal(err)
		return
    }
    defer img.Close()
    w.Header().Set("Content-Type", "image/png")
    io.Copy(w, img)
}
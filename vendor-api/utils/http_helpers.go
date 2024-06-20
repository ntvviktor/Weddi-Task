package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"weddi.org/vendor-api/config"
)

type VendorReviewHelper struct {
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

func GetVendorReviewByReviewId(req *http.Request, reviewId string) ([]byte, error){
	data, err := config.DB.GetVendorReviewByReviewId(req.Context(), reviewId)
	if err != nil {
		log.Fatal("Database Error, Cannot get the resource")
		return nil, err
	}

	var convertedData = VendorReviewHelper{
		ReviewID:     data.ReviewID,
		VendorID:     data.VendorID,
		VendorName:   data.VendorName,
		Poster:       data.Poster,
		Date:         data.Date,
		Rating:       data.Rating,
		Source:       data.Source,
		Content:      data.Content,
		LinkToSource: data.LinkToSource,
	}
	vendorReviews, err := json.Marshal(convertedData)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return vendorReviews, nil
}

func GetVendorReviewByVendorId(req *http.Request, vendorId string) ([]byte, error){
	data, err := config.DB.GetVendorReviewByVendorId(req.Context(), vendorId)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var convertedData []VendorReviewHelper
	for _, vendorReview := range data {
		convertedData = append(convertedData, VendorReviewHelper{
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
		log.Fatal(err)
		return nil, err
	}
	return vendorReviews, nil
}

func ScrapeReviewVendorImage(scrapeUrl string, reviewId string, reviewUrl string) ([]byte, error) {
	contentType := "application/json"
	postBody, _ := json.Marshal(map[string]string{
		"review_id": reviewId,
		"review_url": reviewUrl,
	})

	bodyBuf := bytes.NewBuffer(postBody)
	resp, err := http.Post(scrapeUrl, contentType, bodyBuf)
	if err != nil {
		log.Fatalf("An error occured %v", err)
		return nil, err
	}

	defer resp.Body.Close()
	// Read the response body
	
	body, err := io.ReadAll(resp.Body)
 	if err != nil {
		log.Fatalf("An error occured while reading response body %v", err)
		return nil, err
  	}
	return body, nil
}
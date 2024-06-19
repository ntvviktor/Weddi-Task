package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)


func (apiConfig *apiConfig) getVendorReviewByReviewId(req *http.Request, reviewId string) ([]byte, error){
	data, err := apiConfig.DB.GetVendorReviewByReviewId(req.Context(), reviewId)
	if err != nil {
		apiConfig.logger.Fatal(err)
		return nil, err
	}

	var convertedData = VendorReview{
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
		apiConfig.logger.Fatal(err)
		return nil, err
	}
	return vendorReviews, nil
}

func (apiConfig *apiConfig) getVendorReviewByVendorId(req *http.Request, vendorId string) ([]byte, error){
	data, err := apiConfig.DB.GetVendorReviewByVendorId(req.Context(), vendorId)
	if err != nil {
		apiConfig.logger.Fatal(err)
		return nil, err
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
		return nil, err
	}
	return vendorReviews, nil
}

func scrapeReviewVendorImage(scrapeUrl string, reviewId string, reviewUrl string) ([]byte, error) {
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
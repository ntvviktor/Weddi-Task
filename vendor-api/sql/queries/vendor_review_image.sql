-- name: GetImagesByReviewId :many
SELECT * FROM vendor_review_image
WHERE review_id = ?;
-- name: CreateVendorReview :exec
INSERT INTO vendor_review(review_id, vendor_id, vendor_name, poster, date, rating, source, content, link_to_source)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetAllVendorReviews :many
SELECT * FROM vendor_review;

-- name: GetVendorReviewByVendorId :many
SELECT * FROM vendor_review WHERE vendor_id = ?;

-- name: DeleteVendorReviewByReviewId :exec
DELETE FROM vendor_review WHERE review_id = ?;

-- name: GetVendorReviewByReviewId :one
SELECT * FROM vendor_review WHERE review_id = ? LIMIT 1;

-- name: UpdateVendorReview :exec
UPDATE vendor_review
SET vendor_name = ?,
poster = ?,
date = ?,
rating = ?,
source = ?, 
content = ?,
link_to_source = ?
WHERE review_id = ?;
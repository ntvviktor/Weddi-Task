-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE vendor_review (
    review_id VARCHAR(36) primary key,
    vendor_id VARCHAR(36) not null,
    vendor_name VARCHAR(50) not null, 
    poster VARCHAR(50) not null, 
    date VARCHAR(20) not null, 
    rating INT not null,
    source VARCHAR(50) not null, 
    content VARCHAR(1000) not null, 
    link_to_source VARCHAR(500) not null
);

CREATE TABLE vendor_review_image (
    review_id VARCHAR(36) primary key,
    image_url VARCHAR(500) not null,
    FOREIGN KEY (review_id) REFERENCES vendor_review(review_id)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE vendor_review;
DROP TABLE vendor_review_image;

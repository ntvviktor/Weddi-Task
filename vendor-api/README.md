### CRUD server
A simple RESTful API written in Go to serve data for wedding vendor review.

#### Stack
[sqlc](https://github.com/sqlc-dev/sqlc) compiler to auto generate code that interact with mysql database.
[goose](https://github.com/pressly/goose) to perform database migration in Go.
[Gorilla Mux](https://github.com/gorilla/mux) routing api with net/http compatible library.

#### How to run this project

Pre: 
Create a `.env` file contains the following information:
API_KEY=somekey
DB_CONN=<user>:<password>@/<database-name>
SCRAPE_URL=http://127.0.0.1:5000/api/scrape

1. Clone this project and run:
```
go mod download
```

2. Start the server:
```
go run main.go
```
The server will run on port :8080

package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"weddi.org/vendor-api/internal/database"
)

var DB *database.Queries

func SetupDatabase() *database.Queries {
	_ = godotenv.Load()
	dbUrl := os.Getenv("DB_CONN")
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		fmt.Sprintln("Fail to connect to database")
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	return dbQueries
}
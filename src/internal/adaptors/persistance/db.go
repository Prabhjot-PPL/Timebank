package persistance

import (
	"database/sql"
	"fmt"
	"log"
	"timebank/src/internal/config"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func ConnectToDatabase() (*Database, error) {

	config := config.LoadConfig()

	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", config.DB_User, config.DB_Password, config.DB_Host, config.DB_Port, config.DB_Name)

	db, err := sql.Open("postgres", databaseURL)

	if err != nil {
		return nil, err
	}

	log.Printf("Connected to Database \n\n")

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

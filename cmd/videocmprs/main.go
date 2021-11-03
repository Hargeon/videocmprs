package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Hargeon/videocmprs/api"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/joho/godotenv"
)

const port = ":3001"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	dsn := fmt.Sprintf("user=%s dbname=%s sslmode=%s host=%s port=%s password=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_PASS"))
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}

	h := api.NewHandler(db)
	app := h.InitRoutes()

	if err := app.Listen(port); err != nil {
		log.Fatalln(err)
	}
}

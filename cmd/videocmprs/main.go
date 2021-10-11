package main

import (
	"fmt"
	handler2 "github.com/Hargeon/videocmprs/api"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"os"
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
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	h := handler2.NewHandler(db)
	app := h.InitRoutes()

	if err := app.Listen(port); err != nil {
		log.Fatalln(err)
	}
}

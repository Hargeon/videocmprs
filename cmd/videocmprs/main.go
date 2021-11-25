package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Hargeon/videocmprs/api"
	"github.com/Hargeon/videocmprs/pkg/repository/request"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/Hargeon/videocmprs/pkg/service/broker"
	"github.com/Hargeon/videocmprs/pkg/service/compress"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/joho/godotenv"
)

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

	// init Rabbit publisher
	publisher := broker.NewRabbit(os.Getenv("RABBIT_USER"),
		os.Getenv("RABBIT_PASSWORD"), os.Getenv("RABBIT_HOST"),
		os.Getenv("RABBIT_PORT"))

	publisherConn, err := publisher.Connect("video_convert_test")
	if err != nil {
		log.Fatalln(err)
	}
	defer publisherConn.Close()

	// init Rabbit consumer
	consumer := broker.NewRabbit(os.Getenv("RABBIT_USER"),
		os.Getenv("RABBIT_PASSWORD"), os.Getenv("RABBIT_HOST"),
		os.Getenv("RABBIT_PORT"))

	consumerConn, err := consumer.Connect("video_update_test")
	if err != nil {
		log.Fatalln(err)
	}
	defer consumerConn.Close()

	h := api.NewHandler(db, publisher)
	app := h.InitRoutes()

	msgs, err := consumer.Consume()
	if err != nil {
		log.Fatalln(err)
	}

	reqRepo := request.NewRepository(db)
	vRepo := video.NewRepository(db)
	srv := compress.NewService(reqRepo, vRepo)

	go func() {
		for d := range msgs {
			log.Printf("Received %s", string(d.Body))

			err := srv.UpdateRequest(context.Background(), d.Body)
			log.Println("error occurred when update request", err)

			if err = d.Ack(false); err != nil { // needs to mark a message was processed
				log.Printf("Ack %s", err.Error())
			}
		}
	}()

	if err := app.Listen(os.Getenv("PORT")); err != nil {
		log.Fatalln(err)
	}
}

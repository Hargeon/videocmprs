package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hargeon/videocmprs/api"
	"github.com/Hargeon/videocmprs/pkg/repository/request"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/Hargeon/videocmprs/pkg/service/broker"
	"github.com/Hargeon/videocmprs/pkg/service/cloud"
	"github.com/Hargeon/videocmprs/pkg/service/compress"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"
	"go.uber.org/zap"
)

const migrationsPath = "db/migrations/common"

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln(err)
	}
	defer logger.Sync()

	err = godotenv.Load()
	if err != nil {
		logger.Fatal("godotenv.Load()", zap.String("Error", err.Error()))
	}

	err = runMigrations()
	if err != nil {
		logger.Fatal("error occurred when run migrations", zap.String("Error", err.Error()))
	}

	dsn := os.Getenv("DB_URL")
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		logger.Fatal("", zap.String("Error", err.Error()))
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		logger.Fatal("can't ping db", zap.String("Error", err.Error()))
	}

	// init Rabbit publisher
	publisher := broker.NewRabbit(os.Getenv("RABBIT_URL"))

	publisherConn, err := publisher.Connect("video_convert_test")
	if err != nil {
		logger.Fatal("can't connect to rabbit publisher", zap.String("Error", err.Error()))
	}
	defer publisherConn.Close()

	// init Rabbit consumer
	consumer := broker.NewRabbit(os.Getenv("RABBIT_URL"))

	consumerConn, err := consumer.Connect("video_update_test")
	if err != nil {
		logger.Fatal("can't connect to rabbit consumer", zap.String("Error", err.Error()))
	}
	defer consumerConn.Close()

	msgs, err := consumer.Consume()
	if err != nil {
		logger.Fatal("", zap.String("Error", err.Error()))
	}

	reqRepo := request.NewRepository(db)
	vRepo := video.NewRepository(db)
	srv := compress.NewService(reqRepo, vRepo, logger)

	go func() {
		for d := range msgs {
			logger.Info("Received from rabbit", zap.String("Body", string(d.Body)))

			err := srv.UpdateRequest(context.Background(), d.Body)
			if err != nil {
				logger.Error("Error occurred after updating request status", zap.String("Error", err.Error()))
			}

			if err = d.Ack(false); err != nil { // needs to mark a message was processed
				logger.Error("can't Ack after updating request", zap.String("Error", err.Error()))
			}
		}
	}()

	storage := cloud.NewS3Storage(
		os.Getenv("AWS_BUCKET_NAME"),
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_ACCESS_KEY"),
		os.Getenv("AWS_SECRET_KEY"))

	h := api.NewHandler(db, publisher, storage, logger)
	app := h.InitRoutes()

	logger.Info("Starting web server...")

	go func() {
		if err := app.Listen(os.Getenv("PORT")); err != nil {
			logger.Fatal("Error occurred when starting app", zap.String("Error", err.Error()))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	logger.Info("Server stopped")

	if err := app.Shutdown(); err != nil {
		logger.Fatal("Shutdown server", zap.Error(err))
	}

	logger.Info("Server Exited Properly")
}

func runMigrations() error {
	dsn := os.Getenv("DB_URL")
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return err
	}

	defer db.Close()

	err = goose.Run("up", db, migrationsPath)

	return err
}

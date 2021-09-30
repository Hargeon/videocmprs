package main

import (
	"github.com/Hargeon/videocmprs/pkg/handler"
	"github.com/Hargeon/videocmprs/pkg/service"
	"log"
)

const port = ":3000"

func main() {
	s := service.NewService()
	h := handler.NewHandler(s)
	app := h.InitRoutes()

	if err := app.Listen(port); err != nil {
		log.Fatalln(err)
	}
}

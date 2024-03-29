package server

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	"gitlab.com/metronero/frontend/internal/utils/config"
)

type Server struct {
	*fiber.App
}

func Init() *Server {
	engine := html.New("./views", ".html")

	if config.Debug {
		engine.Reload(true)
		engine.Debug(true)
	}

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	return &Server{app}
}

func (s *Server) StartServerWithGracefulShutdown() {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := s.Shutdown(); err != nil {
			log.Println("Server is not shutting down: ", err)
		}

		close(idleConnsClosed)
	}()

	if err := s.Listen(config.Bind); err != nil {
		log.Println("Failed to start server: ", err)
	}

	<-idleConnsClosed
}

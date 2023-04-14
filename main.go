package main

import (
	"gitlab.com/metronero/frontend/internal/utils/config"
	"gitlab.com/metronero/frontend/internal/utils/server"
)

func main() {
	config.Load()
	app := server.Init()
	app.RegisterRoutes()
	app.StartServerWithGracefulShutdown()
}

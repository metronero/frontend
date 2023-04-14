package main

import (
	"gitlab.com/metronero/frontend/utils/config"
	"gitlab.com/metronero/frontend/utils/server"
)

func main() {
	config.Load()
	app := server.Init()
	app.RegisterRoutes()
	app.StartServerWithGracefulShutdown()
}

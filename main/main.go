package main

import (
	"github.com/spacetimi/passman_server/app_src/app_controller"
	"github.com/spacetimi/passman_server/app_src/app_init"
	"github.com/spacetimi/timi_shared_server/code/core/shared_init"
	"github.com/spacetimi/timi_shared_server/code/server"
)

func main() {

	shared_init.SharedInit(app_init.GetAppInitializer())

	appController := &app_controller.AppController{}
    server.StartServer(appController)
}

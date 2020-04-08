package app_init

import (
    "github.com/spacetimi/passman_server/app_src/app_route_handlers/app_login"
    "github.com/spacetimi/passman_server/app_src/app_route_handlers/home"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
)

type AppController struct { // Implements IAppController
}

func (ac *AppController) RouteHandlers() []controller.IRouteHandler {
    return []controller.IRouteHandler {
        &app_login.AppLoginHandler{},
        &home.HomeHandler{},
    }
}


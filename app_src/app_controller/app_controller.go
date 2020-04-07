package app_controller

import (
    "github.com/spacetimi/timi_shared_server/code/core/controller"
)

type AppController struct { // Implements IAppController
}

func (ac *AppController) RouteHandlers() []controller.IRouteHandler {
    return []controller.IRouteHandler {
        &AppHomeHandler{},
    }
}


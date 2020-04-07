package app_controller

import (
    "fmt"
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "net/http"
)

////////////////////////////////////////////////////////////////////////////////

type AppHomeHandler struct {     // Implements IRouteHandler
}

func (ahh *AppHomeHandler) Routes() []controller.Route {
    return []controller.Route {
        controller.NewRoute(app_routes.Home, []controller.RequestMethodType{controller.GET, controller.POST}),
    }
}

func (ahh *AppHomeHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {
    logger.LogInfo("krisa: here in app home handler")
    _, _ = fmt.Fprintln(httpResponseWriter, "home")
}

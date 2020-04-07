package app_controller

import (
    "fmt"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "net/http"
)

////////////////////////////////////////////////////////////////////////////////

type AppHomeHandler struct {     // Implements IRouteHandler
}

func (ahh *AppHomeHandler) Routes() []controller.Route {
    routes := []controller.Route {
        controller.NewRoute("/home", []controller.RequestMethodType{controller.GET, controller.POST}),
    }
    return routes
}

func (ahh *AppHomeHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {
    logger.LogInfo("krisa: here in app home handler")
    _, _ = fmt.Fprintln(httpResponseWriter, "home")
}

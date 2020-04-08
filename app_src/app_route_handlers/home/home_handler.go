package home

import (
    "fmt"
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "net/http"
)

type HomeHandler struct {     // Implements IRouteHandler
}

func (ahh *HomeHandler) Routes() []controller.Route {
    return []controller.Route {
        controller.NewRoute(app_routes.Home, []controller.RequestMethodType{controller.GET, controller.POST}),
        controller.NewRoute(app_routes.HomeSlash, []controller.RequestMethodType{controller.GET, controller.POST}),
    }
}

func (ahh *HomeHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {
    _, _ = fmt.Fprintln(httpResponseWriter, "home")
}

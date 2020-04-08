package app_login

import (
    "fmt"
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "net/http"
)

type AppLoginHandler struct {     // Implements IRouteHandler
}

func (ahh *AppLoginHandler) Routes() []controller.Route {
    return []controller.Route {
        controller.NewRoute(app_routes.Login, []controller.RequestMethodType{controller.GET, controller.POST}),
    }
}

func (ahh *AppLoginHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {
    _, _ = fmt.Fprintln(httpResponseWriter, "login")
}

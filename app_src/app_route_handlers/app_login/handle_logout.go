package app_login

import (
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/passman_server/app_src/login"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "net/http"
)

func (alh *AppLoginHandler) handleLogout(httpResponseWriter http.ResponseWriter,
                                        request *http.Request,
                                        args *controller.HandlerFuncArgs) {

    login.LogoutUser(httpResponseWriter)
    http.Redirect(httpResponseWriter, request, app_routes.Login, http.StatusSeeOther)
}

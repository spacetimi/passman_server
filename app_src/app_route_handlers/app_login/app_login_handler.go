package app_login

import (
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/passman_server/app_src/app_utils/app_simple_message_page"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "github.com/spacetimi/timi_shared_server/utils/templated_writer"
    "net/http"
    "path/filepath"
)

type AppLoginHandler struct {     // Implements IRouteHandler
    *templated_writer.TemplatedWriter
}

func NewAppLoginHandler() *AppLoginHandler {
    alh := &AppLoginHandler{}
    alh.TemplatedWriter = templated_writer.NewTemplatedWriter(config.GetAppTemplateFilesPath() + "/app_login")

    // Parse templates for every request on LOCAL so that we can iterate over the templates
    // without having to restart the server every time
    alh.TemplatedWriter.ForceReparseTemplates = config.GetEnvironmentConfiguration().AppEnvironment == config.LOCAL

    return alh
}

func (alh *AppLoginHandler) Routes() []controller.Route {
    return []controller.Route {
        controller.NewRoute(app_routes.Login, []controller.RequestMethodType{controller.GET, controller.POST}),
        controller.NewRoute(app_routes.CreateUser, []controller.RequestMethodType{controller.GET, controller.POST}),
        controller.NewRoute(app_routes.ActivateAccount, []controller.RequestMethodType{controller.GET, controller.POST}),
        controller.NewRoute(app_routes.ForgotUsernameOrPassword, []controller.RequestMethodType{controller.GET, controller.POST}),
        controller.NewRoute(app_routes.ResetPassword, []controller.RequestMethodType{controller.GET, controller.POST}),
        controller.NewRoute(app_routes.Logout, []controller.RequestMethodType{controller.GET, controller.POST}),
    }
}

func (alh *AppLoginHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

    /* Handle simple routes */
    switch request.URL.Path {

    case app_routes.Login:
        alh.handleLogin(httpResponseWriter, request, args)
        return

    case app_routes.CreateUser:
        alh.handleCreateUser(httpResponseWriter, request, args)
        return

    case app_routes.ForgotUsernameOrPassword:
        alh.handleForgotUsernameOrPassword(httpResponseWriter, request, args)
        return

    case app_routes.Logout:
        alh.handleLogout(httpResponseWriter, request, args)
        return
    }

    /* Handle routes with path vars */
    if app_routes.ResetPasswordBase == filepath.Dir(request.URL.Path) + "/" {
        alh.handleResetPassword(httpResponseWriter, request, args)
        return
    }
    if app_routes.ActivateAccountBase == filepath.Dir(request.URL.Path) + "/" {
        alh.handleActivateAccount(httpResponseWriter, request, args)
        return
    }

    logger.LogError("unknown route request|request url=" + request.URL.Path)
    httpResponseWriter.WriteHeader(http.StatusNotFound)
}

func showMessage(header string, body string, httpResponseWriter http.ResponseWriter) {
    messageHeader := header
    messageBody := body
    backlinkName := "<< Login"
    app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter,
                                                     messageHeader, messageBody,
                                                     app_routes.Login,
                                                     backlinkName)
}

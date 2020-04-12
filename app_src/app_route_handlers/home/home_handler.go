package home

import (
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/passman_server/app_src/login"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "github.com/spacetimi/timi_shared_server/utils/templated_writer"
    "net/http"
)

type HomeHandler struct {     // Implements IRouteHandler
    *templated_writer.TemplatedWriter
}

func NewHomeHandler() *HomeHandler {
    hh := &HomeHandler{}
    hh.TemplatedWriter = templated_writer.NewTemplatedWriter(config.GetAppTemplateFilesPath() + "/home")

    return hh
}

func (hh *HomeHandler) Routes() []controller.Route {
    return []controller.Route {
        controller.NewRoute(app_routes.Home, []controller.RequestMethodType{controller.GET, controller.POST}),
        controller.NewRoute(app_routes.HomeSlash, []controller.RequestMethodType{controller.GET, controller.POST}),
    }
}

func (hh *HomeHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

    // If user is not logged in, redirect to login page
    user, ok := login.TryGetLoggedInUser(request)
    if !ok {
        http.Redirect(httpResponseWriter, request, app_routes.Login, http.StatusSeeOther)
        return
    }

    pageObject := &HomePageObject{}
    pageObject.Username = user.UserName
    pageObject.UserId = user.UserId

    err := hh.Render(httpResponseWriter,
        "home_page_template.html",
                     pageObject,
                     config.GetEnvironmentConfiguration().AppEnvironment == config.LOCAL)
    if err != nil {
        logger.LogError("error rendering home page template" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}


package app_login

import (
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "github.com/spacetimi/timi_shared_server/utils/templated_writer"
    "net/http"
)

type AppLoginHandler struct {     // Implements IRouteHandler
    *templated_writer.TemplatedWriter
}

func NewAppLoginHandler() *AppLoginHandler {
    alh := &AppLoginHandler{}
    alh.TemplatedWriter = templated_writer.NewTemplatedWriter(config.GetAppTemplateFilesPath() + "/app_login")

    return alh
}

func (alh *AppLoginHandler) Routes() []controller.Route {
    return []controller.Route {
        controller.NewRoute(app_routes.Login, []controller.RequestMethodType{controller.GET, controller.POST}),
    }
}

func (alh *AppLoginHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

    err := alh.TemplatedWriter.Render(httpResponseWriter, "app_login_page_template.html", nil)
    if err != nil {
        logger.LogError("Error executing templates" +
                        "|request url=" + request.URL.String() +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }
}

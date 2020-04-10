package app_login

import (
    "errors"
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "html/template"
    "net/http"
)

var _templates *template.Template

/** Package init **/
func Initialize() error {
    var err error

    _templates, err = template.ParseGlob(config.GetAppTemplateFilesPath() + "/app_login/*")
    if err != nil {
        return errors.New("error parsing templates: " + err.Error())
    }

    return nil
}

////////////////////////////////////////////////////////////////////////////////

type AppLoginHandler struct {     // Implements IRouteHandler
}

func (ahh *AppLoginHandler) Routes() []controller.Route {
    return []controller.Route {
        controller.NewRoute(app_routes.Login, []controller.RequestMethodType{controller.GET, controller.POST}),
    }
}

func (ahh *AppLoginHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

    err := _templates.ExecuteTemplate(httpResponseWriter, "app_login_page_template.html", nil)
    if err != nil {
        logger.LogError("Error executing templates" +
                        "|request url=" + request.URL.String() +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }
}

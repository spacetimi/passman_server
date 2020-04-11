package app_login

import (
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "net/http"
)

func (alh *AppLoginHandler) handleLogin(httpResponseWriter http.ResponseWriter,
                                        request *http.Request,
                                        args *controller.HandlerFuncArgs,
                                        forceReparseTemplates bool) {

    err := alh.TemplatedWriter.Render(httpResponseWriter,
                         "app_login_page_template.html",
                         nil,
                                      forceReparseTemplates)
    if err != nil {
        logger.LogError("Error executing templates" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }
}

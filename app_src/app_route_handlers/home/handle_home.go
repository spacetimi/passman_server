package home

import (
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "net/http"
)

func (hh *HomeHandler) handleHome(user *identity_service.UserBlob, httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {
    pageObject := &HomePageObject{}
    pageObject.Username = user.UserName
    pageObject.UserId = user.UserId

    err := hh.Render(httpResponseWriter,
        "home_page_template.html",
                     pageObject)
    if err != nil {
        logger.LogError("error rendering home page template" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}


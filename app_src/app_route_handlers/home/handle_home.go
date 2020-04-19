package home

import (
    "github.com/spacetimi/passman_server/app_src/data/user_websites"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "net/http"
    "strconv"
)

func (hh *HomeHandler) handleHome(user *identity_service.UserBlob, httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

    userWebsitesBlob, err := user_websites.LoadByUserId(user.UserId, request.Context(), true)
    if err != nil {
        logger.LogError("error getting user websites object" +
                        "|request url=" + request.URL.Path +
                        "|user id=" + strconv.FormatInt(user.UserId, 10) +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }

    pageObject := &HomePageObject{}
    pageObject.Username = user.UserName
    pageObject.UserId = user.UserId

    for _, userWebsite := range userWebsitesBlob.UserWebsites {
        userWebsiteCard := UserWebsiteCardObject{
            WebsiteName:userWebsite.WebsiteName,
        }
        for _, userCredential := range userWebsite.UserWebsiteCredentialsList {
            userWebsiteCard.UserAliases = append(userWebsiteCard.UserAliases, userCredential.UserAlias)
        }

        pageObject.UserWebsiteCards = append(pageObject.UserWebsiteCards, userWebsiteCard)
    }

    err = hh.Render(httpResponseWriter,
       "home_page_template.html",
                    pageObject)
    if err != nil {
        logger.LogError("error rendering home page template" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}


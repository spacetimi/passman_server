package home

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/spacetimi/passman_server/app_src/data/user_secrets"
	"github.com/spacetimi/passman_server/app_src/data/user_websites"
	"github.com/spacetimi/timi_shared_server/code/core/controller"
	"github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"github.com/spacetimi/timi_shared_server/utils/string_utils"
)

func (hh *HomeHandler) handleHome(user *identity_service.UserBlob, httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

	userWebsitesBlob, err := user_websites.LoadByUserId(user.UserId, request.Context(), true)
	if err != nil {
		logger.LogError("error getting user websites blob" +
			"|request url=" + request.URL.Path +
			"|user id=" + strconv.FormatInt(user.UserId, 10) +
			"|error=" + err.Error())
		httpResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	userSecretsBlob, err := user_secrets.LoadByUserId(user.UserId, request.Context(), true)
	if err != nil {
		logger.LogError("error getting user secrets blob" +
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
			WebsiteName:        userWebsite.WebsiteName,
			WebsiteNameEscaped: string_utils.RemoveSpecialCharactersForHtmlId(userWebsite.WebsiteName),
		}
		for _, userCredential := range userWebsite.UserWebsiteCredentialsList {
			userWebsiteCard.UserAliases = append(userWebsiteCard.UserAliases, userCredential.UserAlias)
		}

		pageObject.UserWebsiteCards = append(pageObject.UserWebsiteCards, userWebsiteCard)
	}
	sort.Slice(pageObject.UserWebsiteCards, func(i, j int) bool {
		return pageObject.UserWebsiteCards[i].WebsiteName < pageObject.UserWebsiteCards[j].WebsiteName
	})

	for _, userSecret := range userSecretsBlob.UserSecrets {
		userSecretCard := UserSecretCardObject{
			SecretName:        userSecret.SecretName,
			SecretNameEscaped: string_utils.RemoveSpecialCharactersForHtmlId(userSecret.SecretName),
		}

		pageObject.UserSecretCards = append(pageObject.UserSecretCards, userSecretCard)
	}
	sort.Slice(pageObject.UserSecretCards, func(i, j int) bool {
		return pageObject.UserSecretCards[i].SecretName < pageObject.UserSecretCards[j].SecretName
	})

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

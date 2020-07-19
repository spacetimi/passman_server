package home

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/spacetimi/passman_server/app_src/app_routes"
	"github.com/spacetimi/passman_server/app_src/app_utils/app_simple_message_page"
	"github.com/spacetimi/passman_server/app_src/data/user_websites"
	"github.com/spacetimi/timi_shared_server/code/core/controller"
	"github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

const kDeleteWebsitePostArgWebsiteName = "websiteName"

func (hh *HomeHandler) handleDeleteWebsite(user *identity_service.UserBlob, httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

	parsedArgs, err := parseDeleteWebsitePostArgs(args.PostArgs)
	if err != nil {
		// Show error message and return
		messageHeader := "Something went wrong"
		messageBody := err.Error()
		backlinkName := "<< Home"
		app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
		return
	}

	userWebsites, err := user_websites.LoadByUserId(user.UserId, request.Context(), true)
	if err != nil {
		logger.LogError("error finding user websites blob" +
			"|user id=" + strconv.FormatInt(user.UserId, 10) +
			"|error=" + err.Error())
		// Show error message and return
		messageHeader := "Something went wrong"
		messageBody := "Please try again"
		backlinkName := "<< Home"
		app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
		return
	}

	err = userWebsites.DeleteUserWebsite(parsedArgs.WebsiteName, request.Context())
	if err != nil {
		logger.LogError("error deleting user website from blob" +
			"|user id=" + strconv.FormatInt(user.UserId, 10) +
			"|website name=" + parsedArgs.WebsiteName +
			"|error=" + err.Error())
		// Show error message and return
		messageHeader := "Something went wrong"
		messageBody := "Please try again"
		backlinkName := "<< Home"
		app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
		return
	}

	messageHeader := "Deleted Website: " + parsedArgs.WebsiteName
	messageBody := ""
	backlinkName := "<< Home"
	app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
	return
}

type DeleteWebsitePostArgs struct {
	WebsiteName string
}

func parseDeleteWebsitePostArgs(postArgs map[string]string) (*DeleteWebsitePostArgs, error) {
	websiteName, ok := postArgs[kDeleteUserAliasPostArgWebsiteName]
	if !ok || len(websiteName) == 0 {
		return nil, errors.New("* website name cannot be empty")
	}

	return &DeleteWebsitePostArgs{
		WebsiteName: websiteName,
	}, nil
}

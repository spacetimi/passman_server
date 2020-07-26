package home

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/spacetimi/passman_server/app_src/app_routes"
	"github.com/spacetimi/passman_server/app_src/app_utils/app_simple_message_page"
	"github.com/spacetimi/passman_server/app_src/data/user_files"
	"github.com/spacetimi/timi_shared_server/code/core/controller"
	"github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

const kDeleteFilePostArgFileName = "fileName"

func (hh *HomeHandler) handleDeleteFileFromVault(user *identity_service.UserBlob, httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

	fileName, err := parseDeleteFilePostArgs(args.PostArgs)
	if err != nil {
		// Show error message and return
		messageHeader := "Something went wrong"
		messageBody := err.Error()
		backlinkName := "<< Home"
		app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
		return
	}

	userFilesBlob, err := user_files.LoadByUserId(user.UserId, request.Context(), false)
	if err != nil {
		logger.LogError("error finding user files blob" +
			"|user id=" + strconv.FormatInt(user.UserId, 10) +
			"|error=" + err.Error())
		// Show error message and return
		messageHeader := "Something went wrong"
		messageBody := "Please try again"
		backlinkName := "<< Home"
		app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
		return
	}

	err = userFilesBlob.DeleteFile(fileName, request.Context())
	if err != nil {
		logger.LogError("error deleting user file from vault" +
			"|user id=" + strconv.FormatInt(user.UserId, 10) +
			"|file name=" + fileName +
			"|error=" + err.Error())
		// Show error message and return
		messageHeader := "Something went wrong"
		messageBody := "Please try again"
		backlinkName := "<< Home"
		app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
		return
	}

	// Show success message and return
	messageHeader := "Successfully deleted file from vault: " + fileName
	messageBody := ""
	backlinkName := "<< Home"
	app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
	return
}

func parseDeleteFilePostArgs(postArgs map[string]string) (string, error) {
	fileName, ok := postArgs[kDeleteFilePostArgFileName]
	if !ok || len(fileName) == 0 {
		return "", errors.New("* file name cannot be empty")
	}

	return fileName, nil
}

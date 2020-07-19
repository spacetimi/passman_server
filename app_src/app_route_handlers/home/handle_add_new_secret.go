package home

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/spacetimi/passman_server/app_src/app_routes"
	"github.com/spacetimi/passman_server/app_src/app_utils/app_simple_message_page"
	"github.com/spacetimi/passman_server/app_src/data/user_secrets"
	"github.com/spacetimi/timi_shared_server/code/core/controller"
	"github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

const kAddSecretPostArgSecretName = "secretname"
const kAddSecretPostArgSecretEncryptedValue = "secretvalue"

func (hh *HomeHandler) handleAddNewSecret(user *identity_service.UserBlob, httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

	parsedArgs, err := parseAddNewSecretPostArgs(args.PostArgs)
	if err != nil {
		// Show error message and return
		messageHeader := "Something went wrong"
		messageBody := err.Error()
		backlinkName := "<< Home"
		app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
		return
	}

	userSecretsBlob, err := user_secrets.LoadByUserId(user.UserId, request.Context(), true)
	if err != nil {
		// Show error message and return
		messageHeader := "Something went wrong"
		messageBody := "Please try again"
		backlinkName := "<< Home"
		app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
		return
	}

	err = userSecretsBlob.AddOrModifySecret(parsedArgs.SecretName, parsedArgs.SecretEncrypted, request.Context())
	if err != nil {
		logger.LogError("error add/modify secret" +
			"|user id=" + strconv.FormatInt(user.UserId, 10) +
			"|error=" + err.Error())
		// Show error message and return
		messageHeader := "Something went wrong"
		messageBody := "Please try again"
		backlinkName := "<< Home"
		app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
		return
	}

	// Show Success message and return
	messageHeader := "Success"
	messageBody := "Added new secret: " + parsedArgs.SecretName
	backlinkName := "<< Home"
	app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
	return
}

func parseAddNewSecretPostArgs(postArgs map[string]string) (*AddNewSecretPostArgs, error) {
	secretName, ok := postArgs[kAddSecretPostArgSecretName]
	if !ok || len(secretName) == 0 {
		return nil, errors.New("* secret name cannot be empty")
	}

	secretEncrypted, ok := postArgs[kAddSecretPostArgSecretEncryptedValue]
	if !ok || len(secretEncrypted) == 0 {
		return nil, errors.New("* secret-encrypted cannot be empty")
	}

	return &AddNewSecretPostArgs{
		SecretName:      secretName,
		SecretEncrypted: secretEncrypted,
	}, nil
}

type AddNewSecretPostArgs struct {
	SecretName      string
	SecretEncrypted string
}

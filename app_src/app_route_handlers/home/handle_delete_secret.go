package home

import (
    "errors"
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/passman_server/app_src/app_utils/app_simple_message_page"
    "github.com/spacetimi/passman_server/app_src/data/user_secrets"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "net/http"
    "strconv"
)

const kDeleteSecretPostArgSecretName = "secretname"

func (hh *HomeHandler) handleDeleteSecret(user *identity_service.UserBlob, httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

    secretName, err := parseDeleteSecretPostArgs(args.PostArgs)
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
        logger.LogError("error finding user secrets blob" +
            "|user id=" + strconv.FormatInt(user.UserId, 10) +
            "|error=" + err.Error())
        // Show error message and return
        messageHeader := "Something went wrong"
        messageBody := "Please try again"
        backlinkName := "<< Home"
        app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
        return
    }

    err = userSecretsBlob.DeleteSecret(secretName, request.Context())
    if err != nil {
        logger.LogError("error deleting user secret" +
                        "|user id=" + strconv.FormatInt(user.UserId, 10) +
                        "|secret name=" + secretName +
                        "|error=" + err.Error())
        // Show error message and return
        messageHeader := "Something went wrong"
        messageBody := "Please try again"
        backlinkName := "<< Home"
        app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
        return
    }

    // Show success message and return
    messageHeader := "Successfully deleted secret: " + secretName
    messageBody := ""
    backlinkName := "<< Home"
    app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
    return
}

func parseDeleteSecretPostArgs(postArgs map[string]string) (string, error) {
    secretName, ok := postArgs[kViewSecretPostArgSecretName]
    if !ok || len(secretName) == 0 {
        return "", errors.New("* secret name cannot be empty")
    }

    return secretName, nil
}


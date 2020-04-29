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

const kAddSecretPostArgSecretName = "secretname"
const kAddSecretPostArgSecretValue = "secretvalue"
const kAddSecretPostArgMasterPassword = "masterpassword"

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

    err = userSecretsBlob.AddOrModifySecret(parsedArgs.SecretName, parsedArgs.SecretValue, parsedArgs.MasterPassword, request.Context())
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

    secretValue, ok := postArgs[kAddSecretPostArgSecretValue]
    if !ok || len(secretValue) == 0 {
        return nil, errors.New("* secret cannot be empty")
    }

    masterPassword, ok := postArgs[kAddSecretPostArgMasterPassword]
    if !ok || len(masterPassword) == 0 {
        return nil, errors.New("* master password cannot be empty")
    }

    return &AddNewSecretPostArgs{
                SecretName:secretName,
                SecretValue:secretValue,
                MasterPassword:masterPassword,
            }, nil
}

type AddNewSecretPostArgs struct {
    SecretName string
    SecretValue string
    MasterPassword string
}

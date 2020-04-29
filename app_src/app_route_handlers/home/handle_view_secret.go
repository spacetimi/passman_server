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

const kViewSecretPostArgSecretName = "secretname"
const kViewSecretPostArgMasterPassword = "masterPassword"

func (hh *HomeHandler) handleViewSecret(user *identity_service.UserBlob, httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

    parsedArgs, err := parseViewSecretPostArgs(args.PostArgs)
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

    secretValue, err := userSecretsBlob.GetSecret(parsedArgs.SecretName, parsedArgs.MasterPassword)
    if err != nil {
        logger.LogWarning("error getting secret value" +
                          "|user id=" + strconv.FormatInt(user.UserId, 10) +
                          "|secret name=" + parsedArgs.SecretName +
                          "|error=" + err.Error())
        // Show error message and return
        messageHeader := "Something went wrong"
        messageBody := "Please check your master password and try again"
        backlinkName := "<< Home"
        app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
        return
    }

    pageObject := &ViewSecretPageObject{
        SecretName:parsedArgs.SecretName,
        SecretValue:secretValue,
    }
    err = hh.Render(httpResponseWriter, "view_secret_page_template.html", pageObject)
    if err != nil {
        logger.LogError("error showing view-secret page" +
                        "|user id=" + strconv.FormatInt(user.UserId, 10) +
                        "|secret name=" + parsedArgs.SecretName +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}

func parseViewSecretPostArgs(postArgs map[string]string) (*ViewSecretPostArgs, error) {
    secretName, ok := postArgs[kViewSecretPostArgSecretName]
    if !ok || len(secretName) == 0 {
        return nil, errors.New("* secret name cannot be empty")
    }

    masterPassword, ok := postArgs[kViewSecretPostArgMasterPassword]
    if !ok || len(masterPassword) == 0 {
        return nil, errors.New("* master password cannot be empty")
    }

    return &ViewSecretPostArgs{
                SecretName:secretName,
                MasterPassword:masterPassword,
            }, nil
}

type ViewSecretPostArgs struct {
    SecretName string
    MasterPassword string
}

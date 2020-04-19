package home

import (
    "errors"
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/passman_server/app_src/app_utils/app_simple_message_page"
    "github.com/spacetimi/passman_server/app_src/data/user_websites"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
    "net/http"
)

const kPostArgWebsiteName = "websiteName"
const kPostArgUserAlias   = "userAlias"

func (hh *HomeHandler) handleAddNewWebsite(user *identity_service.UserBlob, httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

    parsedArgs, err := parseAddNewWebsitePostArgs(args.PostArgs)
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
        // Show error message and return
        messageHeader := "Something went wrong"
        messageBody := "Please try again"
        backlinkName := "<< Home"
        app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
        return
    }

    err = userWebsites.AddOrModifyUserWebsiteCredentials(parsedArgs.WebsiteName, parsedArgs.UserAlias, request.Context())
    if err != nil {
        // Show error message and return
        messageHeader := "Something went wrong"
        messageBody := err.Error()
        backlinkName := "<< Home"
        app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
        return
    }

    // Show Success message and return
    messageHeader := "Success"
    messageBody := "Added " + parsedArgs.UserAlias + " @ " + parsedArgs.WebsiteName + " to list of credentials"
    backlinkName := "<< Home"
    app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
    return
}

func parseAddNewWebsitePostArgs(postArgs map[string]string) (*AddNewWebsitePostArgs, error) {
    websiteName, ok := postArgs[kPostArgWebsiteName]
    if !ok || len(websiteName) == 0 {
        return nil, errors.New("* website name cannot be empty")
    }

    userAlias, ok := postArgs[kPostArgUserAlias]
    if !ok || len(userAlias) == 0 {
        return nil, errors.New("* user alias cannot be empty")
    }

    return &AddNewWebsitePostArgs{
                WebsiteName:websiteName,
                UserAlias:userAlias,
            }, nil
}

type AddNewWebsitePostArgs struct {
    WebsiteName string
    UserAlias string
}


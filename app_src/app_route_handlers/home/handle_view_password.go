package home

import (
    "errors"
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/passman_server/app_src/app_utils/app_simple_message_page"
    "github.com/spacetimi/passman_server/app_src/data/user_websites"
    "github.com/spacetimi/passman_server/app_src/password_gen"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "net/http"
    "strconv"
)

const kViewPasswordPostArgWebsiteName = "websiteName"
const kViewPasswordPostArgUserAlias = "userAlias"
const kViewPasswordPostArgMasterPassword = "masterPassword"

const kTemplateName = "view_password_page_template.html"

func (hh *HomeHandler) handleViewPassword(user *identity_service.UserBlob, httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

    parsedArgs, err := parseViewPasswordPostArgs(args.PostArgs)
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

    userWebsite := userWebsites.GetUserWebsite(parsedArgs.WebsiteName)
    if userWebsite == nil {
        // Show error message and return
        logger.LogError("error finding user website object" +
                        "|user id=" + strconv.FormatInt(user.UserId, 10) +
                        "|website name=" + parsedArgs.WebsiteName)
        messageHeader := "Something went wrong"
        messageBody := "Please try again"
        backlinkName := "<< Home"
        app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
        return
    }

    userWebsiteCredentials := userWebsite.GetCredentialsForUserAlias(parsedArgs.UserAlias)
    if userWebsiteCredentials == nil {
        // Show error message and return
        logger.LogError("error finding user website object" +
                        "|user id=" + strconv.FormatInt(user.UserId, 10) +
                        "|website name=" + parsedArgs.WebsiteName +
                        "|user alias=" + parsedArgs.UserAlias)
        messageHeader := "Something went wrong"
        messageBody := "Please try again"
        backlinkName := "<< Home"
        app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
        return
    }

    password, err := password_gen.GeneratePassword(user.UserId,
                                                   user.CreatedTime,
                                                   parsedArgs.WebsiteName,
                                                   parsedArgs.UserAlias,
                                                   userWebsiteCredentials.Version,
                                                   parsedArgs.MasterPassword)
    if err != nil {
        // Show error message and return
        messageHeader := "Something went wrong"
        messageBody := "Please try again"
        backlinkName := "<< Home"
        app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
        return
    }


    // Show password
    pageObject := &ViewPasswordPageObject{}
    pageObject.UserAlias = parsedArgs.UserAlias
    pageObject.WebsiteName = parsedArgs.WebsiteName
    pageObject.Password = password
    err = hh.Render(httpResponseWriter, kTemplateName, pageObject)
    if err != nil {
        logger.LogError("error showing view-password page" +
                        "|user id=" + strconv.FormatInt(user.UserId, 10) +
                        "|website name=" + parsedArgs.WebsiteName +
                        "|user alias=" + parsedArgs.UserAlias +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}

func parseViewPasswordPostArgs(postArgs map[string]string) (*ViewPasswordPostArgs, error) {
    websiteName, ok := postArgs[kViewPasswordPostArgWebsiteName]
    if !ok || len(websiteName) == 0 {
        return nil, errors.New("* website name cannot be empty")
    }

    userAlias, ok := postArgs[kViewPasswordPostArgUserAlias]
    if !ok || len(userAlias) == 0 {
        return nil, errors.New("* user alias cannot be empty")
    }

    masterPassword, ok := postArgs[kViewPasswordPostArgMasterPassword]
    if !ok || len(masterPassword) == 0 {
        return nil, errors.New("* master password cannot be empty")
    }

    return &ViewPasswordPostArgs{
        WebsiteName:websiteName,
        UserAlias:userAlias,
        MasterPassword:masterPassword,
    }, nil
}

type ViewPasswordPostArgs struct {
    WebsiteName string
    UserAlias string
    MasterPassword string
}

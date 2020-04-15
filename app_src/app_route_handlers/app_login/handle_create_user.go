package app_login

import (
    "context"
    "errors"
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/passman_server/app_src/app_utils/app_simple_message_page"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "net/http"
    "regexp"
)

func (alh *AppLoginHandler) handleCreateUser(httpResponseWriter http.ResponseWriter,
                                             request *http.Request,
                                             args *controller.HandlerFuncArgs) {

    pageObject := newCreateUserPageObject()

    if request.Method == controller.POST.String() {
        err := tryCreateNewUser(args.PostArgs, request.Context())
        if err != nil {
            pageObject.SetError(err.Error())
        } else {

            // Show success message and return
            messageHeader := "Successfully created Account"
            messageBody := "Login to continue"
            backlinkName := "<< Login"
            app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter,
                                                             messageHeader, messageBody,
                                                             app_routes.Login,
                                                             backlinkName)
            return
        }
    }

    err := alh.TemplatedWriter.Render(httpResponseWriter,
                         "app_create_user_page_template.html",
                                      pageObject)
    if err != nil {
        logger.LogError("Error executing templates" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }
}

func tryCreateNewUser(postArgs map[string]string, ctx context.Context) error {
    parsed, err := parseCreateUserRequestPostArgs(postArgs)
    if parsed == nil {
        return err
    }

    _, err = identity_service.CreateNewUserByUserNameAndPassword(parsed.Username, parsed.Password, ctx)
    if err != nil {
        return errors.New("* Error creating new user: " + err.Error())
    }

    return nil
}


const kPostArgNewUsername = "new_username"
const kPostArgNewPassword = "new_password"
const kPostArgRetypePassword = "retype_password"

type CreateUserPostArgs struct {
    Username string
    Password string
}

func parseCreateUserRequestPostArgs(postArgs map[string]string) (*CreateUserPostArgs, error) {
    newUsername, ok := postArgs[kPostArgNewUsername]
    if !ok || len(newUsername) == 0 {
        return nil, errors.New("* Please choose a Username")
    }

    newPassword, ok := postArgs[kPostArgNewPassword]
    if !ok || len(newPassword) == 0 {
        return nil, errors.New("* Please choose a Password")
    }

    retypePassword, ok := postArgs[kPostArgRetypePassword]
    if !ok || len(retypePassword) == 0 {
        return nil, errors.New("* Retype Password")
    }

    if newPassword != retypePassword {
        return nil, errors.New("* Passwords do not match")
    }

    passwordErr := validatePassword(newPassword)
    if passwordErr != nil {
        return nil, passwordErr
    }


    return &CreateUserPostArgs{
                Username:newUsername,
                Password:newPassword,
            }, nil
}

func validatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("* Password should be at least 8 characters long")
    }

    regexCheckAlpha := regexp.MustCompile(`[a-zA-Z]+`).MatchString
    if !regexCheckAlpha(password) {
        return errors.New("* Password should contain at least 1 alphabet")
    }

    regexCheckNumeral := regexp.MustCompile(`[0-9]+`).MatchString
    if !regexCheckNumeral(password) {
        return errors.New("* Password should contain at least 1 numeral")
    }

    regexCheckAlphaNumericOnly := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
    if regexCheckAlphaNumericOnly(password) {
        return errors.New("* Password should contain at least 1 special character")
    }

    return nil
}

////////////////////////////////////////////////////////////////////////////////

type CreateUserPageObject struct {
    LoginPageObjectBase
}

func newCreateUserPageObject() *CreateUserPageObject {
    pageObject := &CreateUserPageObject{}
    pageObject.HasError = false
    pageObject.ErrorString = ""

    return pageObject
}

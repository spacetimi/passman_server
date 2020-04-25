package app_login

import (
    "context"
    "errors"
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/passman_server/app_src/app_utils/app_emailer"
    "github.com/spacetimi/passman_server/app_src/app_utils/app_simple_message_page"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
    "github.com/spacetimi/timi_shared_server/utils/email_utils"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "net/http"
)

func (alh *AppLoginHandler) handleForgotUsernameOrPassword(httpResponseWriter http.ResponseWriter,
                                                           request *http.Request,
                                                           args *controller.HandlerFuncArgs) {

    pageObject := &LoginPageObjectBase{}

    if request.Method == controller.POST.String() {
        err := trySendPasswordResetEmail(args.PostArgs, request.Context())
        if err != nil {
            pageObject.SetError(err.Error())

        } else {

            // Show success message and return
            messageHeader := "Sent password-reset instructions to your email address"
            messageBody := ""
            backlinkName := "<< Login"
            app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter,
                                                             messageHeader, messageBody,
                                                             app_routes.Login,
                                                             backlinkName)
            return
        }
    }

    err := alh.TemplatedWriter.Render(httpResponseWriter,
                         "forgot_username_or_password_page_template.html",
                                      pageObject)
    if err != nil {
        logger.LogError("Error executing templates" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }
}

func trySendPasswordResetEmail(postArgs map[string]string, ctx context.Context) error {
    userEmailAddress, err := parsePostArgsForResetPassword(postArgs)
    if err != nil {
        return err
    }

    user, err := identity_service.CheckAndGetUserBlobFromUserEmailAddress(userEmailAddress, ctx)
    if err != nil {
        return errors.New("* Couldn't find any registered user for " + userEmailAddress)
    }

    email := email_utils.Email{
                Subject: "Password reset instructions for your PassMan account",
                Body: "Click here to reset your password.\n" +
                      "This link will be active for 2 days.",
             }

    err = app_emailer.SendEmail(user.UserEmailAddress, email)
    if err != nil {
        return errors.New("* Error sending email. Please try again")
    }

    return nil
}

const kPostArgUserEmailAddress = "userEmailAddress"

func parsePostArgsForResetPassword(postArgs map[string]string) (string, error) {
    userEmailAddress, ok := postArgs[kPostArgUserEmailAddress]
    if !ok || len(userEmailAddress) == 0 {
        return "", errors.New("* Enter email address")
    }

    return userEmailAddress, nil
}

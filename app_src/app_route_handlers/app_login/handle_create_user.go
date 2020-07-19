package app_login

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/spacetimi/passman_server/app_src/app_routes"
	"github.com/spacetimi/passman_server/app_src/app_utils/app_emailer"
	"github.com/spacetimi/passman_server/app_src/app_utils/app_simple_message_page"
	"github.com/spacetimi/passman_server/app_src/login"
	"github.com/spacetimi/timi_shared_server/code/config"
	"github.com/spacetimi/timi_shared_server/code/core/controller"
	"github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
	"github.com/spacetimi/timi_shared_server/utils/email_utils"
	"github.com/spacetimi/timi_shared_server/utils/logger"
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
			messageBody := "Please check your email for the account activation link"
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

	user, err := identity_service.CreateNewUser(parsed.Username, parsed.Password, parsed.EmailAddress, ctx)
	if err != nil {
		return errors.New("* Error creating new user: " + err.Error())
	}

	sendNewAccountActivationEmail(user, ctx)

	return nil
}

func sendNewAccountActivationEmail(user *identity_service.UserBlob, ctx context.Context) {
	newAccountActivationRedisKey, err := login.GenerateNewAccountActivationRedisObject(user, ctx)
	if err != nil {
		logger.LogError("error generating redis object for new account activation" +
			"|user id=" + strconv.FormatInt(user.UserId, 10) +
			"|error=" + err.Error())
		return
	}
	newAccountActivationLink := config.GetEnvironmentConfiguration().ApiServerBaseURL
	if config.GetEnvironmentConfiguration().Port != 80 {
		newAccountActivationLink = newAccountActivationLink + ":" +
			strconv.Itoa(config.GetEnvironmentConfiguration().Port)
	}
	newAccountActivationLink = newAccountActivationLink + app_routes.ActivateAccountBase + newAccountActivationRedisKey

	email := email_utils.Email{
		Subject: "Welcome to PassMan, " + user.UserName + "!",
		Body: "Please click here to activate your account: " + newAccountActivationLink + "\n" +
			// TODO: Do not hard-code 2 days here
			"This link is valid for 2 days",
	}

	err = app_emailer.SendEmail(user.UserEmailAddress, email)
	if err != nil {
		logger.LogError("error sending account activation email" +
			"|user id=" + strconv.FormatInt(user.UserId, 10) +
			"|error=" + err.Error())
	}
}

const kPostArgNewUsername = "new_username"
const kPostArgNewUserEmail = "new_useremail"
const kPostArgNewPassword = "new_password"
const kPostArgRetypePassword = "retype_password"

type CreateUserPostArgs struct {
	Username     string
	Password     string
	EmailAddress string
}

func parseCreateUserRequestPostArgs(postArgs map[string]string) (*CreateUserPostArgs, error) {
	newUsername, ok := postArgs[kPostArgNewUsername]
	if !ok || len(newUsername) == 0 {
		return nil, errors.New("* Please choose a Username")
	}

	newUserEmail, ok := postArgs[kPostArgNewUserEmail]
	if !ok || len(newUserEmail) == 0 {
		return nil, errors.New("* Please enter your email address")
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

	return &CreateUserPostArgs{
		Username:     newUsername,
		Password:     newPassword,
		EmailAddress: newUserEmail,
	}, nil
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

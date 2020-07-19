package app_login

import (
	"context"
	"errors"
	"net/http"

	"github.com/spacetimi/passman_server/app_src/app_routes"
	"github.com/spacetimi/passman_server/app_src/login"
	"github.com/spacetimi/timi_shared_server/code/core/controller"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

const kPostArgUsername = "username"
const kPostArgPassword = "password"

func (alh *AppLoginHandler) handleLogin(httpResponseWriter http.ResponseWriter,
	request *http.Request,
	args *controller.HandlerFuncArgs) {

	pageObject := newLoginPageObject()

	if request.Method == controller.POST.String() {
		err := tryLogin(httpResponseWriter, args.PostArgs, request.Context())
		if err != nil {
			pageObject.SetError(err.Error())
		} else {

			// Redirect to home page
			http.Redirect(httpResponseWriter, request, app_routes.Home, http.StatusSeeOther)
			return
		}
	}

	err := alh.TemplatedWriter.Render(httpResponseWriter,
		"app_login_page_template.html",
		pageObject)
	if err != nil {
		logger.LogError("Error executing templates" +
			"|request url=" + request.URL.Path +
			"|error=" + err.Error())
		httpResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func tryLogin(httpResponseWriter http.ResponseWriter, postArgs map[string]string, ctx context.Context) error {
	parsed, err := parseLoginRequestPostArgs(postArgs)
	if parsed == nil {
		return err
	}

	err = login.TryLoginUserWithCredentials(httpResponseWriter, parsed.Username, parsed.Password, ctx)
	if err != nil {
		return errors.New("wrong username or password")
	}

	return nil
}

func parseLoginRequestPostArgs(postArgs map[string]string) (*LoginPostArgs, error) {
	username, ok := postArgs[kPostArgUsername]
	if !ok || len(username) == 0 {
		return nil, errors.New("* Please enter Username")
	}

	password, ok := postArgs[kPostArgPassword]
	if !ok || len(password) == 0 {
		return nil, errors.New("* Please enter Password")
	}

	return &LoginPostArgs{
		Username: username,
		Password: password,
	}, nil
}

type LoginPostArgs struct {
	Username string
	Password string
}

////////////////////////////////////////////////////////////////////////////////

type LoginPageObject struct {
	LoginPageObjectBase
}

func newLoginPageObject() *LoginPageObject {
	pageObject := &LoginPageObject{}
	pageObject.HasError = false
	pageObject.ErrorString = ""

	return pageObject
}

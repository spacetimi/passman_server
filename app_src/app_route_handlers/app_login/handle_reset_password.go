package app_login

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/spacetimi/passman_server/app_src/app_routes"
	"github.com/spacetimi/passman_server/app_src/app_utils/app_simple_message_page"
	"github.com/spacetimi/passman_server/app_src/data/user_files"
	"github.com/spacetimi/passman_server/app_src/data/user_secrets"
	"github.com/spacetimi/passman_server/app_src/data/user_websites"
	"github.com/spacetimi/passman_server/app_src/login"
	"github.com/spacetimi/timi_shared_server/code/core/adaptors/redis_adaptor"
	"github.com/spacetimi/timi_shared_server/code/core/controller"
	"github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

func (alh *AppLoginHandler) handleResetPassword(httpResponseWriter http.ResponseWriter,
	request *http.Request,
	args *controller.HandlerFuncArgs) {

	pageObject := &LoginPageObjectBase{}

	redisKey, ok := args.RequestPathVars["rediskey"]
	if !ok || len(redisKey) == 0 {
		showMessage("Invalid password-reset link", "", httpResponseWriter)
		return
	}

	userId, err := login.GetUserIdFromResetAccountPasswordRedisKey(redisKey, request.Context())
	if err != nil {
		showMessage("Invalid password-reset link", "", httpResponseWriter)
		return
	}

	user, err := identity_service.GetUserBlobById(userId, request.Context())
	if err != nil {
		showMessage("Unable to fetch user account", "", httpResponseWriter)
		return
	}

	if request.Method == controller.POST.String() {

		err = tryResetPassword(user, args.PostArgs, request.Context())
		if err != nil {
			pageObject.SetError(err.Error())
		} else {

			err = redis_adaptor.Delete(redisKey, request.Context())
			if err != nil {
				logger.LogWarning("error removing password reset link key from redis" +
					"|user id=" + strconv.FormatInt(user.UserId, 10) +
					"|redis key=" + redisKey +
					"|error=" + err.Error())
			}

			// Show success message and return
			messageHeader := "Successfully reset password"
			messageBody := "Please Login again to continue"
			backlinkName := "<< Login"
			app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter,
				messageHeader, messageBody,
				app_routes.Login,
				backlinkName)
			return
		}
	}

	err = alh.TemplatedWriter.Render(httpResponseWriter,
		"reset_password_page_template.html",
		pageObject)
	if err != nil {
		logger.LogError("Error executing templates" +
			"|request url=" + request.URL.Path +
			"|error=" + err.Error())
		httpResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func tryResetPassword(user *identity_service.UserBlob, postArgs map[string]string, ctx context.Context) error {

	password, err := parsePostArgsForResetPassword(postArgs)
	if err != nil {
		return err
	}

	err = identity_service.UpdateUserPassword(user, password, ctx)
	if err != nil {
		return errors.New("* Error updating password. Please try again")
	}

	// Also mark the email address as verified
	// in case the user clicked on reset password immediately after account creation
	// (i.e. without logging in even once)
	err = identity_service.SetUserEmailAddressVerified(user, ctx)
	if err != nil {
		logger.LogError("error marking user email-address as verified" +
			"|user id=" + strconv.FormatInt(user.UserId, 10) +
			"|error=" + err.Error())
	}

	// Also, clear all user data since they cannot be recovered with our client-driven encryption anyway
	clearAllUserData(user.UserId, ctx)

	return nil
}

func clearAllUserData(userId int64, ctx context.Context) {
	{
		userWebsites, err := user_websites.LoadByUserId(userId, ctx, true)
		if err != nil {
			logger.LogError("error fetching user websites blob when trying to reset password" +
				"|user id=" + strconv.FormatInt(userId, 10) +
				"|error=" + err.Error())
		} else {
			err = userWebsites.ClearData(ctx)
			if err != nil {
				logger.LogError("error clearing user websites data on reset password" +
					"|user id=" + strconv.FormatInt(userId, 10) +
					"|error=" + err.Error())
			}
		}
	}

	{
		userSecrets, err := user_secrets.LoadByUserId(userId, ctx, true)
		if err != nil {
			logger.LogError("error fetching user secrets blob when trying to reset password" +
				"|user id=" + strconv.FormatInt(userId, 10) +
				"|error=" + err.Error())
		} else {
			err = userSecrets.ClearData(ctx)
			if err != nil {
				logger.LogError("error clearing user secrets data on reset password" +
					"|user id=" + strconv.FormatInt(userId, 10) +
					"|error=" + err.Error())
			}
		}
	}

	{
		userFiles, err := user_files.LoadByUserId(userId, ctx, true)
		if err != nil {
			logger.LogError("error fetching user files blob when trying to reset password" +
				"|user id=" + strconv.FormatInt(userId, 10) +
				"|error=" + err.Error())
		} else {
			err = userFiles.ClearData(ctx)
			if err != nil {
				logger.LogError("error clearing user files data on reset password" +
					"|user id=" + strconv.FormatInt(userId, 10) +
					"|error=" + err.Error())
			}
		}
	}
}

func parsePostArgsForResetPassword(postArgs map[string]string) (string, error) {
	password, ok := postArgs[kPostArgNewPassword]
	if !ok || len(password) == 0 {
		return "", errors.New("* New Password cannot be empty")
	}

	retypePassword, ok := postArgs[kPostArgRetypePassword]
	if !ok || len(retypePassword) == 0 {
		return "", errors.New("* Retyped-Password cannot be empty")
	}

	if password != retypePassword {
		logger.VarDumpInfo("p", password)
		logger.VarDumpInfo("rp", retypePassword)
		return "", errors.New("* Password and Retyped-Password do not match")
	}

	return password, nil
}

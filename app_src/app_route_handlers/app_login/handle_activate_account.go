package app_login

import (
	"net/http"
	"strconv"

	"github.com/spacetimi/passman_server/app_src/login"
	"github.com/spacetimi/timi_shared_server/code/core/adaptors/redis_adaptor"
	"github.com/spacetimi/timi_shared_server/code/core/controller"
	"github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

func (alh *AppLoginHandler) handleActivateAccount(httpResponseWriter http.ResponseWriter,
	request *http.Request,
	args *controller.HandlerFuncArgs) {

	redisKey, ok := args.RequestPathVars["rediskey"]
	if !ok || len(redisKey) == 0 {
		showMessage("Invalid account-activation link", "", httpResponseWriter)
		return
	}

	userId, err := login.GetUserIdFromNewAccountActivationRedisKey(redisKey, request.Context())
	if err != nil {
		showMessage("Invalid account-activation link", "", httpResponseWriter)
		return
	}

	user, err := identity_service.GetUserBlobById(userId, request.Context())
	if err != nil {
		showMessage("Unable to fetch user account", "", httpResponseWriter)
		return
	}

	err = identity_service.SetUserEmailAddressVerified(user, request.Context())
	if err != nil {
		showMessage("Something went wrong", "Please try again", httpResponseWriter)
		return
	}

	err = redis_adaptor.Delete(redisKey, request.Context())
	if err != nil {
		logger.LogWarning("error removing password reset link key from redis" +
			"|user id=" + strconv.FormatInt(userId, 10) +
			"|redis key=" + redisKey +
			"|error=" + err.Error())
	}

	showMessage("Successfully activated your account", "Please login to continue", httpResponseWriter)
}

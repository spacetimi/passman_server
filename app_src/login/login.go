package login

import (
    "context"
    "errors"
    "github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "net/http"
    "strconv"
    "time"
)

const kCookieName = "passman_login_cookie"
const kSessionExpirationTimeHours = 8

func TryLoginUserWithCredentials(httpResponseWriter http.ResponseWriter, username string, password string, ctx context.Context) error {
    user, err := identity_service.CheckAndGetUserBlobFromUserLoginCredentials(username, password, ctx)
    if err != nil {
        return errors.New("error logging in: " + err.Error())
    }

    userLoginToken, err := identity_service.CreateUserLoginToken(user)
    if err != nil {
        logger.LogError("error creating user login token" +
                        "|user id=" + strconv.FormatInt(user.UserId, 10) +
                        "|error=" + err.Error())
        return errors.New("error creating login token")
    }

    tokenExpiration := time.Now().Add(kSessionExpirationTimeHours * time.Hour)

    cookie := http.Cookie{Name: kCookieName, Value: userLoginToken, Expires: tokenExpiration}
    http.SetCookie(httpResponseWriter, &cookie)

    return nil
}

func TryGetLoggedInUser(request *http.Request) (*identity_service.UserBlob, bool) {
    userLoginToken, err := request.Cookie(kCookieName)
    if err != nil {
        return nil, false
    }

    user, err := identity_service.CheckAndGetUserBlobFromUserLoginToken(userLoginToken.Value, request.Context())
    if err != nil {
        return nil, false
    }

    return user, true
}


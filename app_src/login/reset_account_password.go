package login

import (
    "errors"
    "github.com/spacetimi/passman_server/app_src/app_routes"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/adaptors/redis_adaptor"
    "github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
    "github.com/spacetimi/timi_shared_server/utils/encryption_utils"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "math"
    "math/rand"
    "strconv"
    "time"
)

func GenerateResetAccountPasswordLink(user *identity_service.UserBlob) (string, error) {

    rand.Seed(time.Now().Unix())
    randomString := strconv.FormatInt(user.UserId, 10) + ":" + strconv.Itoa(rand.Intn(math.MaxInt32))
    redisKey := encryption_utils.Generate_md5_hash(randomString)

    err := redis_adaptor.Write(redisKey, strconv.FormatInt(user.UserId, 10), 48 * time.Hour)
    if err != nil {
        logger.LogError("error writing password reset information to redis" +
                        "|user id=" + strconv.FormatInt(user.UserId, 10) +
                        "|error=" + err.Error())
        return "", errors.New("error writing to redis")
    }

    resetPasswordURL := config.GetEnvironmentConfiguration().ApiServerBaseURL + ":" +
                        strconv.Itoa(config.GetEnvironmentConfiguration().Port) +
                        app_routes.ResetPasswordBase + redisKey

    return resetPasswordURL, nil
}

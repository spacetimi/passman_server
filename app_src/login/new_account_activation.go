package login

import (
    "errors"
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

func GenerateNewAccountActivationRedisObject(user *identity_service.UserBlob) (string, error) {

    rand.Seed(time.Now().Unix())
    randomString := strconv.FormatInt(user.UserId, 10) + ":" + strconv.Itoa(rand.Intn(math.MaxInt32))
    redisKey := config.GetAppName() + ":" + "newacct" + ":" + encryption_utils.Generate_md5_hash(randomString)

    err := redis_adaptor.Write(redisKey, strconv.FormatInt(user.UserId, 10), 48 * time.Hour)
    if err != nil {
        logger.LogError("error writing new account activation information to redis" +
                        "|user id=" + strconv.FormatInt(user.UserId, 10) +
                        "|error=" + err.Error())
        return "", errors.New("error writing to redis")
    }

    return redisKey, nil
}

func GetUserIdFromNewAccountActivationRedisKey(redisKey string) (int64, error) {

    redisValue, ok := redis_adaptor.Read(redisKey)
    if !ok || len(redisValue) == 0 {
        return -1, errors.New("error finding redis value for key")
    }

    userId, err := strconv.ParseInt(redisValue, 10, 64)
    if err != nil {
        return -1, errors.New("error parsing redis value into int64")
    }

    return userId, nil
}

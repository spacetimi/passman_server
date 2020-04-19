package password_gen

import (
    "strconv"
)

func GeneratePassword(userId int64, websiteName string, userAlias string, version int, masterPassword string) (string, error) {

    plain := strconv.FormatInt(userId, 10) + websiteName + userAlias + strconv.Itoa(version) + masterPassword
    return plain, nil
}

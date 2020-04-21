package password_gen

import (
    "github.com/spacetimi/timi_shared_server/utils/encryption_utils"
    "github.com/spacetimi/timi_shared_server/utils/string_utils"
    "math/rand"
    "strconv"
)

func GeneratePassword(userId int64, userCreationTime int64, websiteName string, userAlias string, version int, masterPassword string) (string, error) {

    basePasswordBytes := generateBasePasswordBytes(userId, userCreationTime, websiteName, userAlias, version, masterPassword)
    shaHash := encryption_utils.Generate_sha_hash(string(basePasswordBytes))

    // Use the sha-hash as password after making sure it meets common requirements to be used as a password
    seed := userId + userCreationTime + int64(version)
    password := ensurePasswordMeetsRequirements(shaHash, seed)

    return password, nil
}

func generateBasePasswordBytes(userId int64, userCreationTime int64 , websiteName string, userAlias string, version int, masterPassword string) []byte {

    joined := strconv.FormatInt(userId, 10) + strconv.FormatInt(userCreationTime, 10) + websiteName + userAlias + strconv.Itoa(version) + masterPassword
    l := len(joined)

    seed := userId + userCreationTime + int64(version) + addBytesOfString(masterPassword)

    byteArray1 := []byte( string_utils.RandomShuffle(joined, seed) )
    byteArray2 := []byte( string_utils.RandomShuffle(string_utils.Reverse(joined), seed) )

    base := make([]byte, l)
    for i := 0; i < l; i = i + 1 {
        base[i] = byteArray1[i] + byteArray2[i]
    }

    return base
}

const kSpecialCharacters = "!@#~$%^&*()+-={}[]:<>,?"
const kMaxPasswordLength = 10

func ensurePasswordMeetsRequirements(password string, seed int64) string {
    seed = seed + addBytesOfString(password)

    // Generate list of mandatory characters (special characters, upper and lower case alphabets, and numerals)
    mandatoryAppends := ""
    buf := []byte(mandatoryAppends)
    buf = append(buf, generateSpecialCharacter(seed))
    buf = append(buf, generateLowerAlphabet(seed))
    buf = append(buf, generateUpperAlphabet(seed))
    buf = append(buf, generateNumerals(seed))
    mandatoryAppends = string(buf)

    // Rotate the passed in password, and truncate it to ensure final length
    rotated := string_utils.RandomShuffle(password, seed)
    truncatedPassword := string_utils.Truncate(rotated, kMaxPasswordLength - len(mandatoryAppends))

    fullPassword := string_utils.RandomShuffle(truncatedPassword + mandatoryAppends, seed)

    return fullPassword
}

func generateSpecialCharacter(seed int64) byte {
    rand.Seed(seed)
    return kSpecialCharacters[rand.Intn(len(kSpecialCharacters))]
}

func generateLowerAlphabet(seed int64) byte {
    rand.Seed(seed)

    return 'a' + byte(rand.Intn('z' - 'a' + 1))
}

func generateUpperAlphabet(seed int64) byte {
    rand.Seed(seed)

    return 'A' + byte(rand.Intn('Z' - 'A' + 1))
}

func generateNumerals(seed int64) byte {
    rand.Seed(seed)

    return '0' + byte(rand.Intn('9' - '0' + 1))
}

func addBytesOfString(s string) int64 {
    bytes := []byte(s)
    sum := int64(0)
    for _, b := range bytes {
        sum = sum + int64(b)
    }

    return sum
}

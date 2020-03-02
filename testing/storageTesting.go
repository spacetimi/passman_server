package testing

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/spacetimi/timi_shared_server/code/core/services/storage_service"
    "github.com/spacetimi/timi_shared_server/code/core/services/storage_service/storage_typedefs"
    "net/http"
    "strconv"
)

func StorageTestingController(httpResponseWriter http.ResponseWriter, request *http.Request) {
    params := mux.Vars(request)

    if len(params) == 0 {
        storage_testing_home(httpResponseWriter)
        return
    }

    switch params["param1"] {

    case "getuserbooks":
        storage_testing_get_userbooks(httpResponseWriter, params)
        return

    case "getfeatured":
        storage_testing_get(httpResponseWriter, params)
        return

    case "setfeatured":
        storage_testing_set(httpResponseWriter, params)
        return

    case "appendfeatured":
        storage_testing_append_book(httpResponseWriter, params)
        return
    }
}


func storage_testing_home(httpResponseWriter http.ResponseWriter) {
    _, _ = fmt.Fprintln(httpResponseWriter, "storage testing home")
}


func storage_testing_get_userbooks(httpResponseWriter http.ResponseWriter, params map[string]string) {
    userIdString, ok := params["param2"]
    if !ok {
        _, _ = fmt.Fprintln(httpResponseWriter, "no user id specified")
        return
    }

    userId, err := strconv.ParseInt(userIdString, 10, 64)
    if err != nil {
        _, _ = fmt.Fprintln(httpResponseWriter, "error parsing user id: " + err.Error())
        return
    }

    ubb, err := LoadByUserId(userId)
    if err != nil {
        _, _ = fmt.Fprintln(httpResponseWriter, "error loading blob: " + err.Error())
        return
    }

    ubbJson, err := json.MarshalIndent(ubb, "", "  ")
    if err != nil {
        _, _ = fmt.Fprintln(httpResponseWriter, "error serializing blob: " + err.Error())
        return
    }

    _, _ = fmt.Fprintln(httpResponseWriter, "user books blob: " + string(ubbJson))
}

func storage_testing_get(httpResponseWriter http.ResponseWriter, params map[string]string) {
    platformId, ok := params["param2"]
    if !ok {
        _, _ = fmt.Fprintln(httpResponseWriter, "no platform id specified")
        return
    }

    fbb, err := LoadByPlatformId(platformId)
    if err != nil {
        _, _ = fmt.Fprintln(httpResponseWriter, "error getting blob: " + err.Error())
        return
    }

    fbbJson, err := json.MarshalIndent(fbb, "", "  ")
    if err != nil {
        _, _ = fmt.Fprintln(httpResponseWriter, "error serializing blob: " + err.Error())
        return
    }

    _, _ = fmt.Fprintln(httpResponseWriter, "Featured Books Blob: " + string(fbbJson))
}

func storage_testing_set(httpResponseWriter http.ResponseWriter, params map[string]string) {
    platformId, ok := params["param2"]
    if !ok {
        _, _ = fmt.Fprintln(httpResponseWriter, "no platform id specified")
        return
    }

    fbb := NewFeaturedBooksBlob()
    fbb.PlatformId = platformId

    err := storage_service.SetBlob(fbb, context.Background())
    if err != nil {
        _, _ = fmt.Fprintln(httpResponseWriter, "error saving new blob: " + err.Error())
        return
    }

    _, _ = fmt.Fprintln(httpResponseWriter, "successfully inserted new featured books blob")
}

func storage_testing_append_book(httpResponseWriter http.ResponseWriter, params map[string]string) {
    platformId, ok := params["param2"]
    if !ok {
        _, _ = fmt.Fprintln(httpResponseWriter, "no platform id specified")
        return
    }

    bookName, ok := params["param3"]
    if !ok {
        _, _ = fmt.Fprintln(httpResponseWriter, "no book name specified")
        return
    }

    fbb, err := LoadByPlatformId(platformId)
    if err != nil {
        _, _ = fmt.Fprintln(httpResponseWriter, "error loading blob: " + err.Error())
        return
    }

    fbb.FeaturedBooks = append(fbb.FeaturedBooks, FeaturedBook{
        BookName:bookName,
    })
    err = storage_service.SetBlob(fbb, context.Background())
    if err != nil {
        _, _ = fmt.Fprintln(httpResponseWriter, "error setting blob: " + err.Error())
        return
    }

    _, _ = fmt.Fprintln(httpResponseWriter, "successfully appended book to featured books blob")
}

/*** Test Blob ******************************************************************/

type FeaturedBook struct {
    BookName string
    RatingStars int32
}

type FeaturedBooksBlob struct {     // Implements IBlob
    PlatformId string
    FeaturedBooks []FeaturedBook

    storage_typedefs.BlobDescriptor
}

func NewFeaturedBooksBlob() *FeaturedBooksBlob {
    fbb := &FeaturedBooksBlob{}
    fbb.BlobDescriptor = storage_typedefs.NewBlobDescriptor(storage_typedefs.STORAGE_SPACE_APP,
                                                  "featuredbooks",
                                                            []string{"PlatformId"},
                                                            false)
    return fbb
}

func LoadByPlatformId(platformId string) (*FeaturedBooksBlob, error) {
    fbb := NewFeaturedBooksBlob()
    fbb.PlatformId = platformId

    err := storage_service.GetBlobByPrimaryKeys(fbb, context.Background())
    if err != nil {
        return nil, errors.New("error getting blob: " + err.Error())
    }

    return fbb, nil
}


/*** User Identifiable **********************************************************/

type UserIdentifiable struct {
    UserId int64
    storage_typedefs.BlobDescriptor
}

func NewUserIdentifiable(space storage_typedefs.StorageSpace, blobName string, isRedisAllowed bool) UserIdentifiable {
    uib := UserIdentifiable{}
    uib.BlobDescriptor = storage_typedefs.NewBlobDescriptor(space, blobName, []string{"UserId"}, isRedisAllowed)
    return uib
}

/********************************************************************************/

type UserBook struct {
    BookName string
    ChapterId int32
    PageId int32
}

type UserBooksBlob struct {
    UserIdentifiable
    UserBooks []UserBook
}

func NewUserBooksBlob() *UserBooksBlob {
    ubb := &UserBooksBlob{}
    ubb.UserIdentifiable = NewUserIdentifiable(storage_typedefs.STORAGE_SPACE_APP, "userbooks", true)
    return ubb
}

func LoadByUserId(userId int64) (*UserBooksBlob, error) {
    ubb := NewUserBooksBlob()
    ubb.UserId = userId

    err := storage_service.GetBlobByPrimaryKeys(ubb, context.Background())
    if err != nil {
        return nil, errors.New("error getting blob: " + err.Error())
    }

    return ubb, nil
}

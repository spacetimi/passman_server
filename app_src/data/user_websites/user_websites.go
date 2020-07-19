package user_websites

import (
	"context"
	"errors"
	"strconv"

	"github.com/spacetimi/timi_shared_server/code/core/services/storage_service"
	"github.com/spacetimi/timi_shared_server/code/core/services/storage_service/storage_typedefs"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"github.com/spacetimi/timi_shared_server/utils/slice_utils"
)

const kBlobName = "user_websites"
const kVersion = 2

// Implements IBlob
type UserWebsitesBlob struct {
	UserId       int64
	UserWebsites []*UserWebsite

	storage_typedefs.BlobDescriptor `bson:"ignore"`
}

type UserWebsite struct {
	WebsiteName                string
	UserWebsiteCredentialsList []*UserWebsiteCredentials
}

type UserWebsiteCredentials struct {
	UserAlias         string
	PasswordEncrypted string
}

func LoadByUserId(userId int64, ctx context.Context, create bool) (*UserWebsitesBlob, error) {
	userWebsites := newUserWebsitesBlob(userId)

	err := storage_service.GetBlobByPrimaryKeys(userWebsites, ctx)
	if err != nil {
		if !create {
			return nil, errors.New("error getting " + kBlobName + " blob: " + err.Error())
		}

		userWebsites, err = Create(userId, ctx)
		if err != nil {
			return nil, errors.New("error creating " + kBlobName + " blob: " + err.Error())
		}
	}

	return userWebsites, nil
}

func Create(userId int64, ctx context.Context) (*UserWebsitesBlob, error) {
	userWebsites := newUserWebsitesBlob(userId)
	err := storage_service.SetBlob(userWebsites, ctx)
	if err != nil {
		return nil, errors.New("error saving " + kBlobName + " blob: " + err.Error())
	}

	return userWebsites, nil
}

func (blob *UserWebsitesBlob) GetUserWebsite(websiteName string) *UserWebsite {
	for _, websiteCredentials := range blob.UserWebsites {
		if websiteCredentials.WebsiteName == websiteName {
			return websiteCredentials
		}
	}

	return nil
}

func (userWebsite *UserWebsite) GetCredentialsForUserAlias(userAlias string) *UserWebsiteCredentials {
	for _, websiteCredentials := range userWebsite.UserWebsiteCredentialsList {
		if websiteCredentials.UserAlias == userAlias {
			return websiteCredentials
		}
	}

	return nil
}

func (blob *UserWebsitesBlob) AddOrModifyUserWebsiteCredentials(websiteName string, userAlias string, passwordEncrypted string, ctx context.Context) error {
	userWebsite := blob.GetUserWebsite(websiteName)
	if userWebsite == nil {
		userWebsite = &UserWebsite{
			WebsiteName: websiteName,
		}
		blob.UserWebsites = append(blob.UserWebsites, userWebsite)
	}

	credentialsForUserAlias := userWebsite.GetCredentialsForUserAlias(userAlias)
	if credentialsForUserAlias == nil {
		credentialsForUserAlias = &UserWebsiteCredentials{
			UserAlias:         userAlias,
			PasswordEncrypted: passwordEncrypted,
		}
		userWebsite.UserWebsiteCredentialsList = append(userWebsite.UserWebsiteCredentialsList, credentialsForUserAlias)
	} else {
		credentialsForUserAlias.PasswordEncrypted = passwordEncrypted
	}

	// TODO: Avi: Move this somewhere else (like a set-dirty thing for transactions)
	err := storage_service.SetBlob(blob, ctx)
	if err != nil {
		logger.LogError("error saving blob"+
			"|blob name="+kBlobName+
			"|user id="+strconv.FormatInt(blob.UserId, 10),
			"|error="+err.Error())
		return errors.New("error saving changes")
	}

	return nil
}

func (blob *UserWebsitesBlob) DeleteUserWebsiteCredentials(websiteName string, userAlias string, ctx context.Context) error {
	userWebsite := blob.GetUserWebsite(websiteName)
	if userWebsite == nil {
		return errors.New("no such website for user")
	}

	l := len(userWebsite.UserWebsiteCredentialsList)

	index := slice_utils.FindIndexInSlice(l, func(index int) bool {
		return userWebsite.UserWebsiteCredentialsList[index].UserAlias == userAlias
	})
	if index < 0 || index >= l {
		return errors.New("no such user alias")
	}

	userWebsite.UserWebsiteCredentialsList = append(userWebsite.UserWebsiteCredentialsList[:index],
		userWebsite.UserWebsiteCredentialsList[index+1:]...)

	// copy(userWebsite.UserWebsiteCredentialsList[index:], userWebsite.UserWebsiteCredentialsList[index+1:])
	// userWebsite.UserWebsiteCredentialsList = userWebsite.UserWebsiteCredentialsList[:l-1]

	// TODO: Avi: Move this somewhere else (like a set-dirty thing for transactions)
	err := storage_service.SetBlob(blob, ctx)
	if err != nil {
		logger.LogError("error saving blob"+
			"|blob name="+kBlobName+
			"|user id="+strconv.FormatInt(blob.UserId, 10),
			"|error="+err.Error())
		return errors.New("error saving changes")
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func newUserWebsitesBlob(userId int64) *UserWebsitesBlob {
	userWebsites := &UserWebsitesBlob{
		UserId: userId,
	}
	userWebsites.BlobDescriptor = storage_typedefs.NewBlobDescriptor(storage_typedefs.STORAGE_SPACE_APP,
		kBlobName,
		[]string{"UserId"},
		kVersion,
		true)

	return userWebsites
}

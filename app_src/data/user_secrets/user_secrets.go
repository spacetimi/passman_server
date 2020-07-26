package user_secrets

import (
	"context"
	"errors"
	"strconv"

	"github.com/spacetimi/timi_shared_server/code/core/services/storage_service"
	"github.com/spacetimi/timi_shared_server/code/core/services/storage_service/storage_typedefs"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

const kBlobName = "user_secrets"
const kVersion = 1

// Implements IBlob
type UserSecretsBlob struct {
	UserId      int64
	UserSecrets []*UserSecret

	storage_typedefs.BlobDescriptor `bson:"ignore"`
}

type UserSecret struct {
	SecretName          string
	SecretDataEncrypted string
}

func LoadByUserId(userId int64, ctx context.Context, create bool) (*UserSecretsBlob, error) {
	userSecrets := newUserSecretsBlob(userId)

	err := storage_service.GetBlobByPrimaryKeys(userSecrets, ctx)
	if err != nil {
		if !create {
			return nil, errors.New("error getting " + kBlobName + " blob: " + err.Error())
		}

		userSecrets, err = Create(userId, ctx)
		if err != nil {
			return nil, errors.New("error creating " + kBlobName + " blob: " + err.Error())
		}
	}

	return userSecrets, nil

}

func Create(userId int64, ctx context.Context) (*UserSecretsBlob, error) {
	userSecrets := newUserSecretsBlob(userId)
	err := storage_service.SetBlob(userSecrets, ctx)
	if err != nil {
		return nil, errors.New("error saving " + kBlobName + " blob: " + err.Error())
	}

	return userSecrets, nil
}

func (blob *UserSecretsBlob) GetSecret(secretName string) (string, error) {
	userSecret, err := blob.getUserSecretByName(secretName)
	if err != nil {
		return "", err
	}

	return userSecret.SecretDataEncrypted, nil
}

func (blob *UserSecretsBlob) AddOrModifySecret(secretName string, secretDataEncrypted string, ctx context.Context) error {

	userSecret, err := blob.getUserSecretByName(secretName)
	if err != nil || userSecret == nil {
		userSecret = &UserSecret{
			SecretName:          secretName,
			SecretDataEncrypted: secretDataEncrypted,
		}
		blob.UserSecrets = append(blob.UserSecrets, userSecret)
	} else {
		userSecret.SecretDataEncrypted = secretDataEncrypted
	}

	// TODO: Avi: Move this somewhere else (like a set-dirty thing for transactions)
	err = storage_service.SetBlob(blob, ctx)
	if err != nil {
		logger.LogError("error saving blob after add/modify secret" +
			"|blob name=" + kBlobName +
			"|user id=" + strconv.FormatInt(blob.UserId, 10) +
			"|error=" + err.Error())
		return errors.New("error saving changes")
	}

	return nil
}

func (blob *UserSecretsBlob) DeleteSecret(secretName string, ctx context.Context) error {
	var index int
	for i, userSecret := range blob.UserSecrets {
		if userSecret.SecretName == secretName {
			index = i
			break
		}
	}
	if index >= len(blob.UserSecrets) {
		return errors.New("no such secret")
	}

	blob.UserSecrets = append(blob.UserSecrets[:index], blob.UserSecrets[index+1:]...)

	// TODO: Avi: Move this somewhere else (like a set-dirty thing for transactions)
	err := storage_service.SetBlob(blob, ctx)
	if err != nil {
		logger.LogError("error saving blob after deleting secret" +
			"|blob name=" + kBlobName +
			"|user id=" + strconv.FormatInt(blob.UserId, 10) +
			"|error=" + err.Error())
		return errors.New("error saving changes")
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func newUserSecretsBlob(userId int64) *UserSecretsBlob {
	userSecrets := &UserSecretsBlob{
		UserId: userId,
	}
	userSecrets.BlobDescriptor = storage_typedefs.NewBlobDescriptor(storage_typedefs.STORAGE_SPACE_APP,
		kBlobName,
		[]string{"UserId"},
		kVersion,
		true)

	return userSecrets
}

func (blob *UserSecretsBlob) getUserSecretByName(secretName string) (*UserSecret, error) {
	for _, secret := range blob.UserSecrets {
		if secret.SecretName == secretName {
			return secret, nil
		}
	}

	return nil, errors.New("no such secret")
}

package user_files

import (
	"context"
	"errors"
	"strconv"

	"github.com/spacetimi/timi_shared_server/code/core/services/storage_service"
	"github.com/spacetimi/timi_shared_server/code/core/services/storage_service/storage_typedefs"
	"github.com/spacetimi/timi_shared_server/utils/encryption_utils"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

const kBlobName = "user_files"
const kVersion = 1

// Implements IBlob
type UserFilesBlob struct {
	UserId    int64
	UserFiles []*UserFile

	storage_typedefs.BlobDescriptor `bson:"ignore"`
}

type UserFile struct {
	FileName              string
	FileContentsEncrypted string
}

func LoadByUserId(userId int64, ctx context.Context, create bool) (*UserFilesBlob, error) {
	userFiles := newUserFilesBlob(userId)

	err := storage_service.GetBlobByPrimaryKeys(userFiles, ctx)
	if err != nil {
		if !create {
			return nil, errors.New("error getting " + kBlobName + " blob: " + err.Error())
		}

		userFiles, err = Create(userId, ctx)
		if err != nil {
			return nil, errors.New("error creating " + kBlobName + " blob: " + err.Error())
		}
	}

	return userFiles, nil
}

func Create(userId int64, ctx context.Context) (*UserFilesBlob, error) {
	userFiles := newUserFilesBlob(userId)
	err := storage_service.SetBlob(userFiles, ctx)
	if err != nil {
		return nil, errors.New("error saving " + kBlobName + " blob: " + err.Error())
	}

	return userFiles, nil
}

func (blob *UserFilesBlob) AddOrModifyFile(fileName string, fileContents string, filePassword string, ctx context.Context) error {

	fileContentsEncrypted, err := encryption_utils.EncryptUsingAES(fileContents, filePassword)
	if err != nil {
		logger.LogError("error encrypting user-file contents for vault" +
			"|user id=" + strconv.FormatInt(blob.UserId, 10) +
			"|file name=" + fileName +
			"|error=" + err.Error())
		return errors.New("error encrypting file contents")
	}

	userFile, err := blob.getUserFileByName(fileName)
	if err != nil || userFile == nil {
		userFile = &UserFile{
			FileName:              fileName,
			FileContentsEncrypted: fileContentsEncrypted,
		}
		blob.UserFiles = append(blob.UserFiles, userFile)
	} else {
		userFile.FileContentsEncrypted = fileContentsEncrypted
	}

	// TODO: Avi: Move this somewhere else (like a set-dirty thing for transactions)
	err = storage_service.SetBlob(blob, ctx)
	if err != nil {
		logger.LogError("error saving blob after add/modify user-file for vault" +
			"|blob name=" + kBlobName +
			"|user id=" + strconv.FormatInt(blob.UserId, 10) +
			"|error=" + err.Error())
		return errors.New("error saving changes")
	}

	return nil
}

func (blob *UserFilesBlob) DeleteFile(fileName string, ctx context.Context) error {
	var index int
	for i, userFile := range blob.UserFiles {
		if userFile.FileName == fileName {
			index = i
			break
		}
	}
	if index >= len(blob.UserFiles) {
		return errors.New("no such file")
	}

	blob.UserFiles = append(blob.UserFiles[:index], blob.UserFiles[index+1:]...)

	// TODO: Avi: Move this somewhere else (like a set-dirty thing for transactions)
	err := storage_service.SetBlob(blob, ctx)
	if err != nil {
		logger.LogError("error saving blob after deleting file" +
			"|blob name=" + kBlobName +
			"|user id=" + strconv.FormatInt(blob.UserId, 10) +
			"|error=" + err.Error())
		return errors.New("error saving changes")
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func newUserFilesBlob(userId int64) *UserFilesBlob {
	userFiles := &UserFilesBlob{
		UserId: userId,
	}
	userFiles.BlobDescriptor = storage_typedefs.NewBlobDescriptor(storage_typedefs.STORAGE_SPACE_APP,
		kBlobName,
		[]string{"UserId"},
		kVersion,
		true)

	return userFiles
}

func (blob *UserFilesBlob) getUserFileByName(fileName string) (*UserFile, error) {
	for _, file := range blob.UserFiles {
		if file.FileName == fileName {
			return file, nil
		}
	}

	return nil, errors.New("no such file")
}

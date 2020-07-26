package home

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/spacetimi/passman_server/app_src/app_routes"
	"github.com/spacetimi/passman_server/app_src/app_utils/app_simple_message_page"
	"github.com/spacetimi/timi_shared_server/code/core/controller"
	"github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

const kAddFilePostArgFileName = "newFileName"
const kAddFilePostArgFilePassword = "newFilePassword"

func (hh *HomeHandler) handleAddFileToVault(user *identity_service.UserBlob, httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

	parsedArgs, err := parseAddNewFilePostArgs(args.PostArgs)
	if err != nil {
		// Show error message and return
		messageHeader := "Something went wrong"
		messageBody := err.Error()
		backlinkName := "<< Home"
		app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
		return
	}
	logger.VarDumpInfo("file name", parsedArgs.FileName)

	var buffer bytes.Buffer
	uploadedFile, _, err := request.FormFile("newFile")
	if err != nil {
		logger.LogError("error getting file from request for adding to vault" +
			"|error=" + err.Error())
		httpResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		err = uploadedFile.Close()
	}()

	_, err = io.Copy(&buffer, uploadedFile)
	if err != nil {
		logger.LogError("error copying file contents from request for adding to vault" +
			"|error=" + err.Error())
		httpResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileContents := buffer.String()
	logger.VarDumpInfo("uploaded file contents", fileContents)

	http.Redirect(httpResponseWriter, request, "/", http.StatusOK)
}

func parseAddNewFilePostArgs(postArgs map[string]string) (*AddNewFilePostArgs, error) {

	fileName, ok := postArgs[kAddFilePostArgFileName]
	if !ok || len(fileName) == 0 {
		return nil, errors.New("* file name cannot be empty")
	}

	filePassword, ok := postArgs[kAddFilePostArgFilePassword]
	if !ok || len(filePassword) == 0 {
		return nil, errors.New("* file password cannot be empty")
	}

	return &AddNewFilePostArgs{
		FileName:     fileName,
		FilePassword: filePassword,
	}, nil
}

type AddNewFilePostArgs struct {
	FileName     string
	FilePassword string
}

package home

import (
	"bytes"
	"io"
	"net/http"

	"github.com/spacetimi/timi_shared_server/code/core/controller"
	"github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

func (hh *HomeHandler) handleAddFileToVault(user *identity_service.UserBlob, httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		logger.LogError("error parsing request for adding file to vault" +
			"|error=" + err.Error())
		// TODO: Add file name to logger here
		httpResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	var buffer bytes.Buffer
	uploadedFile, _, err := request.FormFile("fileForVault")
	if err != nil {
		logger.LogError("error getting file from request for adding to vault" +
			"|error=" + err.Error())
		// TODO: Add file name to logger here
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
		// TODO: Add file name to logger here
		httpResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileContents := buffer.String()
	logger.VarDumpInfo("uploaded file contents", fileContents)

	http.Redirect(httpResponseWriter, request, "/", http.StatusOK)
}

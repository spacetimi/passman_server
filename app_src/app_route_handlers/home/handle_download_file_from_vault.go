package home

import (
	"bytes"
	"errors"
	"net/http"
	"time"

	"github.com/spacetimi/passman_server/app_src/app_routes"
	"github.com/spacetimi/passman_server/app_src/app_utils/app_simple_message_page"
	"github.com/spacetimi/passman_server/app_src/data/user_files"
	"github.com/spacetimi/timi_shared_server/code/core/controller"
	"github.com/spacetimi/timi_shared_server/code/core/services/identity_service"
)

const kDownloadFilePostArgFileName = "fileName"
const kDownloadFilePostArgFilePassword = "filePassword"

func (hh *HomeHandler) handleDownloadFileFromVault(user *identity_service.UserBlob, httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

	parsedArgs, err := parseDownloadFilePostArgs(args.PostArgs)
	if err != nil {
		// Show error message and return
		messageHeader := "Something went wrong"
		messageBody := err.Error()
		backlinkName := "<< Home"
		app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
		return
	}

	userFilesBlob, err := user_files.LoadByUserId(user.UserId, request.Context(), true)
	if err != nil {
		// Show error message and return
		messageHeader := "Something went wrong"
		messageBody := "Please try again"
		backlinkName := "<< Home"
		app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
		return
	}

	userFileContents, err := userFilesBlob.GetUserFileContentsByName(parsedArgs.FileName, parsedArgs.FilePassword)
	if err != nil {
		// Show error message and return
		messageHeader := "No such file: " + parsedArgs.FileName
		messageBody := ""
		backlinkName := "<< Home"
		app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
		return
	}

	// Mark the returned content as downloadable to the browser
	httpResponseWriter.Header().Add("Content-Disposition", "Attachment; filename="+parsedArgs.FileName)

	http.ServeContent(httpResponseWriter, request, parsedArgs.FileName, time.Now(), bytes.NewReader(userFileContents))
}

func parseDownloadFilePostArgs(postArgs map[string]string) (*DownloadFilePostArgs, error) {

	fileName, ok := postArgs[kDownloadFilePostArgFileName]
	if !ok || len(fileName) == 0 {
		return nil, errors.New("* file name cannot be empty")
	}

	filePassword, ok := postArgs[kDownloadFilePostArgFilePassword]
	if !ok || len(filePassword) == 0 {
		return nil, errors.New("* file password cannot be empty")
	}

	return &DownloadFilePostArgs{
		FileName:     fileName,
		FilePassword: filePassword,
	}, nil
}

type DownloadFilePostArgs struct {
	FileName     string
	FilePassword string
}

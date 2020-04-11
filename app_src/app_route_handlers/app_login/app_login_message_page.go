package app_login

import (
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "github.com/spacetimi/timi_shared_server/utils/templated_writer"
    "net/http"
)

func ShowAppLoginMessagePage(httpResponseWriter http.ResponseWriter,
                             messageHeader string,
                             messageBody string,
                             backlinkHref string,
                             backlinkHrefName string) {

    pageObject := &AppLoginMessagePageObject{
        MessageHeader:messageHeader,
        MessageBody:messageBody,
        BackLinkHref:backlinkHref,
        BackLinkHrefName:backlinkHrefName,
    }

    page := newAppLoginMessagePage()
    err := page.Render(httpResponseWriter,
          "app_login_message_page_template.html",
                       pageObject,
                       config.GetEnvironmentConfiguration().AppEnvironment == config.LOCAL)
    if err != nil {
        logger.LogError("error showing login message page|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}

////////////////////////////////////////////////////////////////////////////////

type AppLoginMessagePage struct {
    *templated_writer.TemplatedWriter
}

func newAppLoginMessagePage() *AppLoginMessagePage {
    almp := &AppLoginMessagePage{}
    almp.TemplatedWriter = templated_writer.NewTemplatedWriter(config.GetAppTemplateFilesPath() + "/app_login")

    return almp
}


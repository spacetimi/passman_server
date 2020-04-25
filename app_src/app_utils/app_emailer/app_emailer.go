package app_emailer

import "github.com/spacetimi/timi_shared_server/utils/email_utils"

var _instance *email_utils.Emailer

func Initialize() {
    // TODO: Move email account credentials to AWS secrets-manager
    account := email_utils.EmailAccount{
        EmailAddress:"support-mailer@spacetimi.com",
        Password:"spacetimi1!",
    }
    _instance = email_utils.NewEmailer(account, "smtp.gmail.com", 587)
}

func SendEmail(to string, email email_utils.Email) error {
    return _instance.SendEmail(to, email)
}


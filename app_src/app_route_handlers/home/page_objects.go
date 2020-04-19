package home

type HomePageObject struct {
    HasError bool
    ErrorString string

    Username string
    UserId int64

    UserWebsiteCards []UserWebsiteCardObject
}

type UserWebsiteCardObject struct {
    WebsiteName string
    UserAliases []string
}

func (po *HomePageObject) SetError(errorString string) {
    po.HasError = true
    po.ErrorString = errorString
}

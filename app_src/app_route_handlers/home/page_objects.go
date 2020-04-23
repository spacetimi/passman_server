package home

type PageObjectBase struct {
    HasError bool
    ErrorString string

    Username string
    UserId int64
}

type HomePageObject struct {
    PageObjectBase

    UserWebsiteCards []UserWebsiteCardObject
}

type UserWebsiteCardObject struct {
    WebsiteName string
    UserAliases []string
}

type ViewPasswordPageObject struct {
    PageObjectBase

    UserAlias string
    WebsiteName string
    Password string
}

func (po *PageObjectBase) SetError(errorString string) {
    po.HasError = true
    po.ErrorString = errorString
}

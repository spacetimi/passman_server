package home

type PageObjectBase struct {
	HasError    bool
	ErrorString string

	Username string
	UserId   int64
}

type HomePageObject struct {
	PageObjectBase

	UserWebsiteCards []UserWebsiteCardObject
	UserSecretCards  []UserSecretCardObject
}

type UserWebsiteCardObject struct {
	WebsiteName        string
	WebsiteNameEscaped string
	UserAliases        []UserAliasCardObject
}

type UserAliasCardObject struct {
	Alias        string
	AliasEscaped string
}

type UserSecretCardObject struct {
	SecretNameEscaped string
	SecretName        string
}

type ViewPasswordPageObject struct {
	PageObjectBase

	UserAlias         string
	WebsiteName       string
	PasswordEncrypted string
}

type ViewSecretPageObject struct {
	PageObjectBase

	SecretName      string
	SecretEncrypted string
}

func (po *PageObjectBase) SetError(errorString string) {
	po.HasError = true
	po.ErrorString = errorString
}

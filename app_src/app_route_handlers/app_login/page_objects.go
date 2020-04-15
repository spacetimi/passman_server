package app_login

type LoginPageObjectBase struct {
    HasError bool
    ErrorString string
}

func (po *LoginPageObjectBase) SetError(errorString string) {
    po.HasError = true
    po.ErrorString = errorString
}



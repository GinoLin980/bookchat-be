package internalerror

import "errors"

var ErrUserHashError = errors.New("can't hash user password")
var ErrPasswordInvalid = errors.New("wrong password")
var ErrUserForbidden = errors.New("forbidden")

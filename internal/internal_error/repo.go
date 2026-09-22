package internalerror

import "errors"

var ErrRecordNotFound = errors.New("record not found")
var ErrDatabaseErr = errors.New("database error")
var ErrDuplicatedError = errors.New("duplicated key")

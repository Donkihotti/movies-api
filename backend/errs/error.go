package errs 

import (
	"errors"
)

var NotFound = errors.New("entity not found")
var BadRequest = errors.New("error: bad request")
var ServerError = errors.New("internal server error")
var Conflict = errors.New("error: conflict")

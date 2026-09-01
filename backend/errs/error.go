package errs

import (
	"errors"
)

//types of possible errors
var NotFound = errors.New("entity not found")
var BadRequest = errors.New("error: bad request")
var ServerError = errors.New("internal server error")
var Conflict = errors.New("error: conflict")


//error struct which will get the detailed message what went wrong
type ErrorStruct struct {
	ErrType error
	ErrMsg 	error 
}



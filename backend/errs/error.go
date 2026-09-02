package errs

import (
	"errors"
	"fmt"
)

//types of possible errors
var NotFound = errors.New("Not Found")
var BadRequest = errors.New("Bad request")
var ServerError = errors.New("Server error")
var Conflict = errors.New("Conflict")


//error struct which will get the detailed message what went wrong
type ErrorStruct struct {
	ErrType error //type of the error above
	ErrMsg 	error 
}

func NewErrorStruct(errType, errMsg error) (newStruct ErrorStruct) {
	newStruct.ErrType = errType
	newStruct.ErrMsg = errMsg 
	return newStruct

}

func (es *ErrorStruct) Error() string {
	return fmt.Sprintf("Error: %v\nMessage: %v", es.ErrType, es.ErrMsg)
}

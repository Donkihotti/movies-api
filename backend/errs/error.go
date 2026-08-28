package errs 

import (
	"errors"
)
//HERE COMES THE GLOBAL ERROR HANDLING 


var ErrNotFound = errors.New("entity not found")
var ErrInvalidInput = errors.New("incorrect input")
var ErrServerError = errors.New("internal server error")


func GlobalErrorHandler(err error) (error, string) {
	
	switch errors.Is(err) {
		case ErrNotFound: return err, 404
		case ErrInvalidInput: return err, 400
		case ErrServerError: return err, 500
		
	}

}



package errors

// ErrorHandlerService can provide the HTTP status code appropriate
// for a given error.
type ErrorHandlerService interface {
	GetHTTPStatusCode(error) int
}

// NewErrorHandlerService creates an ErrorHandlerService.
func NewErrorHandlerService() ErrorHandlerService {
	return errorHandlerService{}
}

type errorHandlerService struct{}

func (ehs errorHandlerService) GetHTTPStatusCode(err error) int {
	switch err.(type) {
	case harkBadRequestError:
		return 400
	case harkNotFoundError:
		return 404
	case harkConflictError:
		return 409
	case harkInternalServerError:
		return 500
	default:
		return 500
	}
}

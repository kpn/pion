package s3responses

import "net/http"

type Error struct {
	Code       string
	Message    string
	Key        string
	BucketName string
	Resource   string
	RequestID  string `xml:"RequestId" json:"RequestId"`
}

var InternalErrorResponse = Error{
	Code:    "InternalError",
	Message: "We encountered an internal error. Please try again.",
}

// MapHttpCodes translates from s3 errorCode to http status code
var MapHttpCodes = map[string]int{
	"BucketAlreadyExists": http.StatusConflict,
	"InternalError":       http.StatusInternalServerError,
	// add other errorCodes here
}

// NewError creates a S3 Error struct object based on the errorCode parameter
func NewError(errorCode string, params ...interface{}) *Error {
	switch errorCode {
	case "BucketAlreadyExists":
		bucketName := params[0].(string)
		return &Error{
			Code:       errorCode,
			Message:    "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.",
			BucketName: bucketName,
			RequestID:  "", // TODO add tracking request ID
		}
		// add other errorCode use-case here
	}
	return &InternalErrorResponse
}

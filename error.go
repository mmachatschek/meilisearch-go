package meilisearch

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// ErrCode are all possible errors found during requests
type ErrCode int

const (
	// ErrCodeUnknown default error code, undefined
	ErrCodeUnknown ErrCode = 0
	// ErrCodeMarshalRequest impossible to serialize request body
	ErrCodeMarshalRequest ErrCode = iota + 1
	// ErrCodeResponseUnmarshalBody impossible deserialize the response body
	ErrCodeResponseUnmarshalBody
	// MeilisearchApiError send by the Meilisearch api
	MeilisearchApiError
	// MeilisearchApiError send by the Meilisearch api
	MeilisearchApiErrorWithoutMessage
	// MeilisearchTimeoutError
	MeilisearchTimeoutError
	// MeilisearchCommunicationError impossible execute a request
	MeilisearchCommunicationError
)

const (
	rawStringCtx                               = `(path "${method} ${endpoint}" with method "${function}")`
	rawStringMarshalRequest                    = `unable to marshal body from request: '${request}'`
	rawStringResponseUnmarshalBody             = `unable to unmarshal body from response: '${response}' status code: ${statusCode}`
	rawStringMeilisearchApiError               = `unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, MeilisearchApiError Message: ${message}, ErrorCode: ${errorCode}, ErrorType: ${errorType}, ErrorLink: ${errorLink}`
	rawStringMeilisearchApiErrorWithoutMessage = `unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, MeilisearchApiError Message: ${message}`
	rawStringMeilisearchTimeoutError           = `MeilisearchTimeoutError`
	rawStringMeilisearchCommunicationError     = `MeilisearchCommunicationError unable to execute request`
)

func (e ErrCode) rawMessage() string {
	switch e {
	case ErrCodeMarshalRequest:
		return rawStringMarshalRequest + " " + rawStringCtx
	case ErrCodeResponseUnmarshalBody:
		return rawStringResponseUnmarshalBody + " " + rawStringCtx
	case MeilisearchApiError:
		return rawStringMeilisearchApiError + " " + rawStringCtx
	case MeilisearchApiErrorWithoutMessage:
			return rawStringMeilisearchApiErrorWithoutMessage + " " + rawStringCtx
	case MeilisearchTimeoutError:
		return rawStringMeilisearchTimeoutError + " " + rawStringCtx
	case MeilisearchCommunicationError:
		return rawStringMeilisearchCommunicationError + " " + rawStringCtx
	default:
		return rawStringCtx
	}
}

type meilisearchApiMessage struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
	ErrorType string `json:"errorType"`
	ErrorLink string `json:"errorLink"`
}

// Error is the internal error structure that all exposed method use.
// So ALL errors returned by this library can be cast to this struct (as a pointer)
type Error struct {
	// Endpoint is the path of the request (host is not in)
	Endpoint string

	// Method is the HTTP verb of the request
	Method string

	// Function name used
	Function string

	// RequestToString is the raw request into string ('empty request' if not present)
	RequestToString string

	// RequestToString is the raw request into string ('empty response' if not present)
	ResponseToString string

	// Error info from Meilisearch api
	// Message is the raw request into string ('empty meilisearch message' if not present)
	MeilisearchApiMessage meilisearchApiMessage

	// StatusCode of the request
	StatusCode int

	// StatusCode expected by the endpoint to be considered as a success
	StatusCodeExpected []int

	rawMessage string

	// OriginError is the origin error that produce the current Error. It can be nil in case of a bad status code.
	OriginError error

	// ErrCode is the internal error code that represent the different step when executing a request that can produce
	// an error.
	ErrCode ErrCode
}

// Error return a well human formatted message.
func (e Error) Error() string {
	message := namedSprintf(e.rawMessage, map[string]interface{}{
		"endpoint":           e.Endpoint,
		"method":             e.Method,
		"function":           e.Function,
		"request":            e.RequestToString,
		"response":           e.ResponseToString,
		"statusCodeExpected": e.StatusCodeExpected,
		"statusCode":         e.StatusCode,
		"message":            e.MeilisearchApiMessage.Message,
		"errorCode":          e.MeilisearchApiMessage.ErrorCode,
		"errorType":          e.MeilisearchApiMessage.ErrorType,
		"errorLink":          e.MeilisearchApiMessage.ErrorLink,
	})
	if e.OriginError != nil {
		return errors.Wrap(e.OriginError, message).Error()
	}

	return message
}

// WithErrCode add an error code to an error
func (e *Error) WithErrCode(err ErrCode, errs ...error) *Error {
	if errs != nil {
		e.OriginError = errs[0]
	}

	e.rawMessage = err.rawMessage()
	e.ErrCode = err
	return e
}

// ErrorBody add a body to an error
func (e *Error) ErrorBody(body []byte) {
	e.ResponseToString = string(body)
	msg := meilisearchApiMessage{}
	err := json.Unmarshal(body, &msg)
	if err == nil {
		e.MeilisearchApiMessage.Message = msg.Message
		e.MeilisearchApiMessage.ErrorCode = msg.ErrorCode
		e.MeilisearchApiMessage.ErrorType = msg.ErrorType
		e.MeilisearchApiMessage.ErrorLink = msg.ErrorLink
	}
}

func namedSprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		format = strings.ReplaceAll(format, "${"+key+"}", fmt.Sprintf("%v", val))
	}
	return format
}

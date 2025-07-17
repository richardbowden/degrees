package problems

// very close to RFC 7807 Error handling

import (
	"fmt"
	"net/http"
)

type Kind int

const (
	Other          Kind = iota // Unclassified error. This value is not printed in the error message.
	Invalid                    // Invalid operation for this type of item.
	IO                         // External I/O error such as network failure.
	Exist                      // Item already exists.
	NotExist                   // Item does not exist.
	Private                    // Information withheld.
	Internal                   // Internal error or inconsistency.
	BrokenLink                 // Link target does not exist.
	Database                   // Error from database.
	Validation                 // Input validation error.
	Unanticipated              // Unanticipated error.
	InvalidRequest             // Invalid Request

	// Unauthenticated is used when a request lacks valid authentication credentials.
	//
	// For Unauthenticated errors, the response body will be empty.
	// The error is logged and http.StatusUnauthorized (401) is sent.
	Unauthenticated // Unauthenticated Request

	// Unauthorized is used when a user is authenticated, but is not authorized
	// to access the resource.
	//
	// For Unauthorized errors, the response body should be empty.
	// The error is logged and http.StatusForbidden (403) is sent.
	Unauthorized
)

func (k Kind) String() string {
	switch k {
	case Other:
		return "other error"
	case Invalid:
		return "invalid operation"
	case IO:
		return "I/O error"
	case Exist:
		return "item already exists"
	case NotExist:
		return "item does not exist"
	case BrokenLink:
		return "link target does not exist"
	case Private:
		return "information withheld"
	case Internal:
		return "internal error"
	case Database:
		return "database error"
	case Validation:
		return "input validation error"
	case Unanticipated:
		return "unanticipated error"
	case InvalidRequest:
		return "invalid request error"
	case Unauthenticated:
		return "unauthenticated request"
	case Unauthorized:
		return "unauthorized request"
	}
	return "unknown error kind"
}

func httpErrorStatusCode(k Kind) int {
	// the zero value of Kind is Other, so if no Kind is present
	// in the error, Other is used. Errors should always have a
	// Kind set, otherwise, a 500 will be returned and no
	// error message will be sent to the caller
	switch k {
	case NotExist:
		return http.StatusNotFound
	case Invalid, Exist, Private, BrokenLink, Validation, InvalidRequest:
		return http.StatusBadRequest
	case Other, IO, Internal, Database, Unanticipated:
		return http.StatusInternalServerError
	case Unauthenticated:
		return http.StatusUnauthorized
	case Unauthorized:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

type Title string

type Problem struct {
	Status int      `json:"status"`
	Title  Title    `json:"title" example:"Bad Request" doc:"A short, human-readable summary of the problem type. This value should not change between occurrences of the error."`
	Detail string   `json:"detail"`
	Errors []Detail `json:"errors,omitempty"`

	Kind    Kind  `json:"-"`
	OrigErr error `json:"-"`
}

func (p *Problem) AddDetail(err error) {
	if converted, ok := err.(Detail); ok {
		p.Errors = append(p.Errors, converted)
		return
	}
	p.Errors = append(p.Errors, Detail{Message: err.Error()})
}

type Detail struct {
	Message  string `json:"message"`
	Value    string `json:"value,omitempty"`
	Location string `json:"-"`
}

func (d Detail) Error() string {
	if d.Value == "" {
		return d.Message
	}
	return fmt.Sprintf("%s (%s)", d.Message, d.Value)
}

func New(kind Kind, msg string, errs ...error) Problem {
	status := httpErrorStatusCode(kind)
	p := Problem{
		Status: status,
		Title:  Title(http.StatusText(status)),
		Kind:   kind,
		Detail: msg,
	}

	for _, e := range errs {
		p.AddDetail(e)
	}

	return p
}

func (p Problem) Error() string {
	return p.Detail
}

func (p Problem) GetStatus() int {
	return p.Status
}

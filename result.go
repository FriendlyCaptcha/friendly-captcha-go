package friendlycaptcha

import (
	"errors"
	"fmt"
)

// VerifyResult wraps the response from the Friendly Captcha API when verifying a captcha, making it easier
// to work with. In the simplest case, you can just check `ShouldAccept` to see if the captcha was solved correctly or
// if you should accept it anyway (e.g. because the API was down).
type VerifyResult struct {
	// Success is true if the captcha was solved correctly.
	Success bool

	// Status is the HTTP Response status code of the request to the Friendly Captcha API.
	Status int

	response VerifyResponse
	strict   bool

	// The error that occurred during verification, if any.
	err error
}

// RequestError returns the error, if any (nil otherwise).
func (r VerifyResult) RequestError() error {
	return r.err
}

// Strict returns whether the verification was strict.
//
// If strict is false (= the default), and verification was not able to happen (e.g. because your API key is incorrect, or the Friendly Captcha API is down)
// then `ShouldAccept` will return true regardless.
func (r VerifyResult) Strict() bool {
	return r.strict
}

// ShouldAccept returns true if you should allow the request to pass through.
// It is possible that verification wasn't possible, perhaps the API is unavailable. In that case this function will
// also return true, unless you enable `strict` mode for the client.
func (r VerifyResult) ShouldAccept() bool {
	if r.WasAbleToVerify() {
		return r.response.Success
	}
	if r.err != nil {
		if r.strict { // If Strict mode is enabled, we do not accept any captcha if there was an error.
			return false
		}
		if errors.Is(r.err, ErrVerificationRequest) ||
			errors.Is(
				r.err,
				ErrVerificationFailedDueToClientError,
			) { // Failure to talk to Friendly Captcha verification API or client error (e.g. wrong API key)
			return true
		}
		return false
	}

	panic(
		"Implementation error in friendly-captcha-go ShouldAccept: error should never be nil if success is false. " + fmt.Sprintf(
			"%+v",
			r,
		),
	)
}

// ShouldReject is the inverse of ShouldAccept.
func (r VerifyResult) ShouldReject() bool {
	return !r.ShouldAccept()
}

// IsRequestError returns true if an error occurred while sending the request to the Friendly Captcha API or interpreting its response.
// This could be due to network connectivity issues, or the Friendly Captcha API experiencing downtime.
func (r VerifyResult) IsRequestError() bool {
	return r.err != nil && errors.Is(r.err, ErrVerificationRequest)
}

// This is an error that is not due to a connection error, but due to a client error (e.g. wrong API key).
// You should log this and notify yourself and fix this as soon as possible.
//
// It's usually still a good idea to accept the captcha: it's better to accept any captcha than to lock all users out.
func (r VerifyResult) IsErrorDueToClientError() bool {
	return r.err != nil && errors.Is(r.err, ErrVerificationFailedDueToClientError)
}

// Response returns the response from the Friendly Captcha API.
func (r VerifyResult) Response() VerifyResponse {
	return r.response
}

// HTTPStatusCode returns the HTTP status code of the response from the Friendly Captcha API.
func (r VerifyResult) HTTPStatusCode() int {
	return r.Status
}

// WasAbleToVerify returns true if the captcha could be verified. If this is false, you should log the reason why
// and investigate (you can retrieve the error using the `RequestError` method). The `IsErrorDueToClientError` method
// will tell you if the error was due to a client error (e.g. wrong API key) - which will require your action to fix.
func (r VerifyResult) WasAbleToVerify() bool {
	return r.Status == 200 && !r.IsRequestError()
}

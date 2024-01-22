package friendlycaptcha

import (
	"errors"
	"fmt"
)

type VerifyResult struct {
	Success bool

	// HTTP Response status code
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

func (r VerifyResult) ShouldAccept() bool {
	if r.WasAbleToVerify() {
		return r.response.Success
	}
	if r.err != nil {
		if r.strict { // If Strict mode is enabled, we do not accept any captcha if there was an error.
			return false
		}
		if errors.Is(r.err, ErrVerificationRequest) || errors.Is(r.err, ErrVerificationFailedDueToClientError) { // Failure to talk to Friendly Captcha verification API or client error (e.g. wrong API key)
			return true
		}
		return false
	}

	panic("Implementation error in friendly-captcha-go-sdk ShouldAccept: error should never be nil if success is false. " + fmt.Sprintf("%+v", r))
}

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

func (r VerifyResult) Response() VerifyResponse {
	return r.response
}

func (r VerifyResult) HTTPStatusCode() int {
	return r.Status
}

func (r VerifyResult) WasAbleToVerify() bool {
	return r.Status == 200 && !r.IsRequestError()
}

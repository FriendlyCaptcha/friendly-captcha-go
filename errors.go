package friendlycaptcha

import "errors"

// Could not create the request body (i.e. JSON marshal it), this should never happen but if it does then probably
// the captcha solution value is really weird - let's not accept the verification.
var ErrCreatingVerificationRequest = errors.New("could not create verification request body")

// The POST request to the Friendly Captcha API could not be completed for some reason.
var ErrVerificationRequest = errors.New("verification request failed talking to Friendly Captcha API")

// This error signifies a non-200 response from the server. Usually this means that your API key was wrong.
// You should notify yourself if this happens, but it's usually still a good idea to accept the captcha even though
// we were unable to verify it: we don't want to lock users out.
var ErrVerificationFailedDueToClientError = errors.New("verification request failed due to a client error (check your credentials)")

type ErrorCode string

const (
	// (401) You forgot to set the X-API-Key header.
	ErrorCodeAuthRequired ErrorCode = "auth_required"
	// (401) The API key you provided is invalid.
	ErrorCodeAuthInvalid ErrorCode = "auth_invalid"
	// (400) The sitekey in your request is invalid.
	ErrorCodeSitekeyInvalid ErrorCode = "sitekey_invalid"
	// (400) The response field is missing in your request.
	ErrorCodeResponseMissing ErrorCode = "response_missing"

	// (200) The response field is invalid.
	ErrorCodeResponseInvalid ErrorCode = "response_invalid"
	// (200) The response has expired.
	ErrorCodeResponseTimeout ErrorCode = "response_timeout"
	// (200) The response has already been used.
	ErrorCodeResponseDuplicate ErrorCode = "response_duplicate"

	// (400) Something else is wrong with your request, e.g. the request body was empty.
	ErrorCodeBadRequest ErrorCode = "bad_request"
)

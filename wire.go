package friendlycaptcha

import "time"

// VerifyRequest is the request body for the /api/v2/captcha/siteverify endpoint. As a user of the SDK
// you generally don't need to create this struct yourself, instead you should use the Client's methods.
type VerifyRequest struct {
	// The response value that the user submitted in the frc-captcha-response field.
	Response string `json:"response"`
	// Optional: the sitekey that you want to make sure the puzzle was generated from.
	Sitekey string `json:"sitekey,omitempty"`
}

// VerifyResponseChallengeData is the data found in the challenge field of a VerifyResponse.
// It contains information about the challenge that was solved.
type VerifyResponseChallengeData struct {
	Timestamp time.Time `json:"timestamp"`
	Origin    string    `json:"origin"`
}

// VerifyResponseData is the data found in the data field of a VerifyResponse.
type VerifyResponseData struct {
	Challenge VerifyResponseChallengeData `json:"challenge"`
}

// VerifyResponseError is the data found in the error field of a VerifyResponse in case of an error.
type VerifyResponseError struct {
	ErrorCode ErrorCode `json:"error_code"`
	Detail    string    `json:"detail"`
}

// VerifyResponse is the response body for the /api/v2/captcha/siteverify endpoint. This is what the Friendly
// Captcha API returns.
type VerifyResponse struct {
	Success bool `json:"success"`

	// This field is only present when the success field is true.
	Data *VerifyResponseData `json:"data,omitempty"`
	// This field is only present when the success field is false.
	Error *VerifyResponseError `json:"error,omitempty"`
}

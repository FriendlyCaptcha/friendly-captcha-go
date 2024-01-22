package friendlycaptcha

import "time"

type VerifyRequest struct {
	// The response value that the user submitted in the frc-captcha-response field.
	Response string `json:"response"`
	// Optional: the sitekey that you want to make sure the puzzle was generated from.
	Sitekey string `json:"sitekey,omitempty"`
}

type VerifyResponseChallengeData struct {
	Timestamp time.Time `json:"timestamp"`
	Origin    string    `json:"origin"`
}

type VerifyResponseData struct {
	Challenge VerifyResponseChallengeData `json:"challenge"`
}

type VerifyResponseError struct {
	ErrorCode ErrorCode `json:"error_code"`
	Detail    string    `json:"detail"`
}

type VerifyResponse struct {
	Success bool `json:"success"`

	// This field is only present when the success field is true.
	Data *VerifyResponseData `json:"data,omitempty"`
	// This field is only present when the success field is false.
	Error *VerifyResponseError `json:"error,omitempty"`
}

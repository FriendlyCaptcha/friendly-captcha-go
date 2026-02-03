package friendlycaptcha

import (
	"encoding/json"
	"time"

	"github.com/guregu/null/v6"
)

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
	// EventID is a unique identifier for this siteverify call.
	EventID string `json:"event_id"`

	// Challenge contains information about the challenge that was solved.
	Challenge VerifyResponseChallengeData `json:"challenge"`

	// RiskIntelligenceRaw contains risk information about the solver of the captcha.
	// This may be `null` if risk intelligence is not enabled for your Friendly Captcha account.
	//
	// Note this is the raw JSON data, you probably want to use the RiskIntelligence field instead. This field is
	// available in case you need to access fields that are not yet modeled in the SDK.
	RiskIntelligenceRaw null.Value[json.RawMessage] `json:"risk_intelligence"`

	// RiskIntelligence contains risk information about the solver of the captcha.
	// This may be `null` if risk intelligence is not enabled for your Friendly Captcha account.
	RiskIntelligence null.Value[RiskIntelligenceData] `json:"-"`
}

// UnmarshalJSON implements custom JSON unmarshaling for VerifyResponseData.
// It automatically populates the RiskIntelligence field from RiskIntelligenceRaw.
func (v *VerifyResponseData) UnmarshalJSON(data []byte) error {
	// Use an auxiliary struct to avoid infinite recursion
	type Alias VerifyResponseData
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(v),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Populate RiskIntelligence from RiskIntelligenceRaw
	if v.RiskIntelligenceRaw.Valid {
		var riskData RiskIntelligenceData
		if err := json.Unmarshal(v.RiskIntelligenceRaw.V, &riskData); err != nil {
			return err
		}
		v.RiskIntelligence = null.ValueFrom(riskData)
	}

	return nil
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

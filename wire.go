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

// RiskIntelligenceRetrieveRequest is the request body for the /api/v2/riskIntelligence/retrieve endpoint.
type RiskIntelligenceRetrieveRequest struct {
	Token string `json:"token"`
}

// RiskIntelligenceTokenData is metadata about the risk intelligence token in a retrieve response.
type RiskIntelligenceTokenData struct {
	// Timestamp when the token was generated.
	Timestamp time.Time `json:"timestamp"`
	// Timestamp when the token expires.
	ExpiresAt time.Time `json:"expires_at"`
	// Number of times the token has been used.
	NumUses int64 `json:"num_uses"`
	// The origin of the site where the token was generated.
	Origin string `json:"origin"`
}

// RiskIntelligenceRetrieveResponseData is the data field in a successful retrieve response.
type RiskIntelligenceRetrieveResponseData struct {
	// EventID is a unique identifier for this risk intelligence retrieve call.
	EventID string `json:"event_id"`

	// Token contains metadata about the token used for retrieval.
	Token RiskIntelligenceTokenData `json:"token"`

	// RiskIntelligenceRaw contains the raw JSON risk information extracted from the provided token.
	//
	// Note this is the raw JSON data, you probably want to use the RiskIntelligence field instead. This field is
	// available in case you need to access fields that are not yet modeled in the SDK.
	RiskIntelligenceRaw null.Value[json.RawMessage] `json:"risk_intelligence"`

	// RiskIntelligence contains risk information extracted from the provided token.
	RiskIntelligence null.Value[RiskIntelligenceData] `json:"-"`
}

// UnmarshalJSON implements custom JSON unmarshaling for VerifyResponseData.
// It automatically populates the RiskIntelligence field from RiskIntelligenceRaw.
func (r *RiskIntelligenceRetrieveResponseData) UnmarshalJSON(data []byte) error {
	// Use an auxiliary struct to avoid infinite recursion
	type Alias RiskIntelligenceRetrieveResponseData
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Populate RiskIntelligence from RiskIntelligenceRaw
	if r.RiskIntelligenceRaw.Valid {
		var riskData RiskIntelligenceData
		if err := json.Unmarshal(r.RiskIntelligenceRaw.V, &riskData); err != nil {
			return err
		}
		r.RiskIntelligence = null.ValueFrom(riskData)
	}

	return nil
}

// RiskIntelligenceRetrieveResponse is the response body for the /api/v2/riskIntelligence/retrieve endpoint.
type RiskIntelligenceRetrieveResponse struct {
	Success bool `json:"success"`

	// This field is only present when the success field is true.
	Data *RiskIntelligenceRetrieveResponseData `json:"data,omitempty"`
	// This field is only present when the success field is false.
	Error *VerifyResponseError `json:"error,omitempty"`
}

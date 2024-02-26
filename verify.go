package friendlycaptcha

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// VerifyCaptchaResponse takes a captcha response and verifies it with the Friendly Captcha API.
// It returns a VerifyResult, which contains the result of the verification.
//
// On this result struct that will allow you to check if the verification could be performed, and whether
// you should allow the user to proceed.
func (frc *Client) VerifyCaptchaResponse(ctx context.Context, captchaResponse string) VerifyResult {
	result := VerifyResult{}
	reqBody := VerifyRequest{
		Response: captchaResponse,
		Sitekey:  frc.Sitekey,
	}
	result.strict = frc.Strict
	// We should never end up with this status code, unless we fail to be able to marshal the request body.
	result.Status = -1

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		result.err = fmt.Errorf("%w: %v", ErrCreatingVerificationRequest, err)
		return result
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, frc.SiteverifyEndpoint, bytes.NewReader(reqBodyJSON))
	if err != nil {
		result.err = fmt.Errorf("%w: %v", ErrCreatingVerificationRequest, err)
		return result
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", frc.APIKey)
	req.Header.Set("X-Frc-Sdk", fmt.Sprintf("friendly-captcha-go-sdk@%s", Version))

	resp, err := frc.HTTPClient.Do(req)
	if err != nil {
		result.err = fmt.Errorf("%w: %v", ErrVerificationRequest, err)
		return result
	}
	defer resp.Body.Close()
	result.Status = resp.StatusCode

	decoder := json.NewDecoder(resp.Body)
	var vr VerifyResponse
	err = decoder.Decode(&vr)
	if err != nil {
		// This is likely a problem with the Friendly Captcha API.
		// It's returning something unexpected.
		result.err = fmt.Errorf("%w: %v", ErrVerificationRequest, err)
		return result
	}

	if resp.StatusCode != http.StatusOK {
		// Intentionally let this through, it's probably a problem in our credentials
		result.err = fmt.Errorf("%w [status %d]: %s", ErrVerificationFailedDueToClientError, resp.StatusCode, vr.Error)
		return result
	}

	result.response = vr
	result.Success = vr.Success
	return result
}

package friendlycaptcha

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// A ClientOption is a function that can be passed to NewClient to configure a new Client.
type ClientOption func(*Client) error

// A client for the Friendly Captcha API, see also the API docs at https://developer.friendlycaptcha.com
type Client struct {
	APIKey      string
	Sitekey     string
	APIEndpoint string
	// If Strict is set to true only strictly verified captcha response will be allowed.
	// For example: if your server can not reach the Friendly Captcha endpoint, it will still advise to accept the response
	// regardless.
	//
	// By default Strict is false: `ShouldAccept()` will return true when for instance the Friendly Captcha API
	// could not be reached.
	Strict bool
	// The HTTP client to use for making requests to the Friendly Captcha API.
	// Defaults to `http.DefaultClient`
	HTTPClient *http.Client
}

// The name of the form field that, by default, the widget will put the captcha response in.
const ResponseFormFieldName = "frc-captcha-response"

const (
	globalAPIEndpoint = "https://global.frcapi.com"
	euAPIEndpoint     = "https://eu.frcapi.com"
)

var (
	errCreateRequest = errors.New("failed to create HTTP request")
)

// NewClient creates a new Friendly Captcha client with the given options.
func NewClient(opts ...ClientOption) (*Client, error) {
	const (
		defaultAPIEndpoint = globalAPIEndpoint
	)

	c := &Client{
		HTTPClient:  http.DefaultClient,
		APIEndpoint: defaultAPIEndpoint,
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated *Client as the argument
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	if c.APIKey == "" {
		return nil, fmt.Errorf(
			"you must set your Friendly Captcha API key using `WithAPIKey()` when creating a new client",
		)
	}

	return c, nil
}

// WithAPIKey sets the API key for the client.
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.APIKey = apiKey
		return nil
	}
}

// WithSitekey sets the sitekey for the client. This is optional.
func WithSitekey(sitekey string) ClientOption {
	return func(c *Client) error {
		c.Sitekey = sitekey
		return nil
	}
}

// In strict mode only strictly verified captcha response are allowed. If your API key is invalid or your server can not reach the API endpoint all requests will be rejected.
//
// This defaults to `false`.
func WithStrictMode(strict bool) ClientOption {
	return func(c *Client) error {
		c.Strict = strict
		return nil
	}
}

// WithAPIEndpoint sets the API endpoint domain for the client.
// Takes a domain without path (e.g., "https://global.frcapi.com"), or the shorthands "global" or "eu".
func WithAPIEndpoint(apiEndpoint string) ClientOption {
	return func(c *Client) error {
		switch apiEndpoint {
		case "global":
			apiEndpoint = globalAPIEndpoint
		case "eu":
			apiEndpoint = euAPIEndpoint
		case "":
			return fmt.Errorf("apiEndpoint must not be empty")
		}
		c.APIEndpoint = apiEndpoint
		return nil
	}
}

// WithSiteverifyEndpoint sets the API endpoint for the client.
// Deprecated: Use WithAPIEndpoint instead. This function strips the path from the URL and calls WithAPIEndpoint.
// Takes a full URL, or the shorthands "global" or "eu".
func WithSiteverifyEndpoint(siteverifyEndpoint string) ClientOption {
	return func(c *Client) error {
		if siteverifyEndpoint == "" {
			return fmt.Errorf("siteverifyEndpoint must not be empty")
		}

		// Handle shorthands
		if siteverifyEndpoint == "global" || siteverifyEndpoint == "eu" {
			return WithAPIEndpoint(siteverifyEndpoint)(c)
		}

		// Parse URL to extract scheme and host (domain without path)
		u, err := url.Parse(siteverifyEndpoint)
		if err != nil {
			return fmt.Errorf("invalid siteverifyEndpoint URL: %w", err)
		}

		// Construct the API endpoint without path
		apiEndpoint := u.Scheme + "://" + u.Host

		return WithAPIEndpoint(apiEndpoint)(c)
	}
}

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
	result.Status = -1

	var vr VerifyResponse
	statusCode, err := frc.postJSON(ctx, "/api/v2/captcha/siteverify", reqBody, &vr)
	result.Status = statusCode
	if err != nil {
		if errors.Is(err, errCreateRequest) {
			result.err = fmt.Errorf("%w: %v", ErrCreatingVerificationRequest, err)
			return result
		}
		result.err = fmt.Errorf("%w: %v", ErrVerificationRequest, err)
		return result
	}

	if statusCode != http.StatusOK {
		// Intentionally let this through, it's probably a problem in our credentials
		result.err = fmt.Errorf("%w [status %d]: %s", ErrVerificationFailedDueToClientError, statusCode, vr.Error)
		return result
	}

	result.response = vr
	result.Success = vr.Success
	return result
}

// RetrieveRiskIntelligence takes a risk intelligence token and retrieves the associated risk intelligence data from the Friendly Captcha API.
// It returns a RiskIntelligenceRetrieveResult, which contains the risk intelligence data.
func (frc *Client) RetrieveRiskIntelligence(ctx context.Context, token string) RiskIntelligenceRetrieveResult {
	result := RiskIntelligenceRetrieveResult{}
	reqBody := RiskIntelligenceRetrieveRequest{
		Token:   token,
		Sitekey: frc.Sitekey,
	}
	// We should never end up with this status code, unless we fail to be able to marshal the request body.
	result.Status = -1

	var retrieveResponse RiskIntelligenceRetrieveResponse
	statusCode, err := frc.postJSON(ctx, "/api/v2/riskIntelligence/retrieve", reqBody, &retrieveResponse)
	result.Status = statusCode
	if err != nil {
		if errors.Is(err, errCreateRequest) {
			result.err = fmt.Errorf("%w: %v", ErrCreatingRiskIntelligenceRetrieveRequest, err)
			return result
		}
		result.err = fmt.Errorf("%w: %v", ErrRiskIntelligenceRetrieveRequest, err)
		return result
	}

	if statusCode != http.StatusOK {
		// Intentionally let this through, it's probably a problem in our credentials.
		result.err = fmt.Errorf(
			"%w [status %d]: %+v",
			ErrRiskIntelligenceRetrieveFailedDueToClientError,
			statusCode,
			retrieveResponse.Error,
		)
		return result
	}

	result.response = retrieveResponse
	result.Success = retrieveResponse.Success
	return result
}

func (frc *Client) postJSON(ctx context.Context, path string, requestBody any, responseBody any) (int, error) {
	reqBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return -1, fmt.Errorf("%w: %v", errCreateRequest, err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		frc.APIEndpoint+path,
		bytes.NewReader(reqBodyJSON),
	)
	if err != nil {
		return -1, fmt.Errorf("%w: %v", errCreateRequest, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", frc.APIKey)
	req.Header.Set("Frc-Sdk", fmt.Sprintf("friendly-captcha-go@%s", Version))

	resp, err := frc.HTTPClient.Do(req)
	if err != nil {
		return -1, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return statusCode, fmt.Errorf("error decoding response body: %v", err)
	}

	return statusCode, nil
}

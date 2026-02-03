package friendlycaptcha

import (
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

package friendlycaptcha

import (
	"fmt"
	"net/http"
)

type ClientOption func(*Client) error

// A client for the Friendly Captcha API, see also the API docs at https://developer.friendlycaptcha.com
type Client struct {
	APIKey             string
	Sitekey            string
	SiteverifyEndpoint string
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
	globalSiteverifyEndpointURL = "https://global.frcapi.com/api/v2/captcha/siteverify"
	euSiteverifyEndpointURL     = "https://eu.frcapi.com/api/v2/captcha/siteverify"
)

func NewClient(opts ...ClientOption) (*Client, error) {
	const (
		defaultSiteverifyEndpoint = globalSiteverifyEndpointURL
	)

	c := &Client{
		HTTPClient:         http.DefaultClient,
		SiteverifyEndpoint: defaultSiteverifyEndpoint,
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
		return nil, fmt.Errorf("you must set your Friendly Captcha API key using `WithAPIKey()` when creating a new client")
	}

	return c, nil
}

func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.APIKey = apiKey
		return nil
	}
}

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

// Takes a full URL, or the shorthands `"global"` or `"eu"` .
func WithSiteverifyEndpoint(siteverifyEndpoint string) ClientOption {
	if siteverifyEndpoint == "global" {
		siteverifyEndpoint = globalSiteverifyEndpointURL
	} else if siteverifyEndpoint == "eu" {
		siteverifyEndpoint = euSiteverifyEndpointURL
	}

	return func(c *Client) error {
		c.SiteverifyEndpoint = siteverifyEndpoint
		return nil
	}
}

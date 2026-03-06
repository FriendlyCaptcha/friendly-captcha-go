//go:build sdkintegration

package friendlycaptcha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	MockServerURL                             = "http://localhost:1090"
	CaptchaSiteverifyTestCasesJSONEndpoint    = "/api/v1/captcha/siteverifyTests"
	RiskIntelligenceRetrieveTestCasesEndpoint = "/api/v1/riskIntelligence/retrieveTests"
)

var runTestSDKMockServerTests = true

type CaptchaSiteverifyTestCasesFile struct {
	Version int                         `json:"version"`
	Tests   []CaptchaSiteverifyTestCase `json:"tests"`
}

type CaptchaSiteverifyTestCase struct {
	Name        string `json:"name"`
	Response    string `json:"response"`
	Expectation struct {
		ShouldAccept    bool `json:"should_accept"`
		WasAbleToVerify bool `json:"was_able_to_verify"`
		IsClientError   bool `json:"is_client_error"`
	} `json:"expectation"`
	Strict             bool            `json:"strict"`
	SiteverifyResponse json.RawMessage `json:"siteverify_response"`
}

type SuccessCaptchaSiteverifyResponse struct {
	Data VerifyResponseData `json:"data"`
}

type RiskIntelligenceRetrieveTestCasesFile struct {
	Version int                                `json:"version"`
	Tests   []RiskIntelligenceRetrieveTestCase `json:"tests"`
}

type RiskIntelligenceRetrieveTestCase struct {
	Name        string `json:"name"`
	Token       string `json:"token"`
	Expectation struct {
		WasAbleToRetrieve bool `json:"was_able_to_retrieve"`
		IsClientError     bool `json:"is_client_error"`
	} `json:"expectation"`
	RetrieveResponse json.RawMessage `json:"retrieve_response"`
}

type SuccessRiskIntelligenceRetrieveResponse struct {
	Data RiskIntelligenceRetrieveResponseData `json:"data"`
}

func loadTestCasesFromServer[T any](testCasesEndpoint string) (T, error) {
	var zero T

	resp, err := http.Get(MockServerURL + testCasesEndpoint)
	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) {
			return zero, fmt.Errorf("sdk mock testing server is not running, please run it first")
		}
		return zero, fmt.Errorf("failed to load test cases from mock server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("failed to load test cases from mock server: %s", resp.Status)
	}

	var testCases T
	err = json.NewDecoder(resp.Body).Decode(&testCases)
	if err != nil {
		return zero, fmt.Errorf("failed to decode test cases from mock server: %w", err)
	}

	return testCases, nil
}

func TestSDKWithMockServerCaptchaSiteverify(t *testing.T) {
	if !runTestSDKMockServerTests {
		t.Skip("Skipping SDK tests, run the mock server first and set runTestSDKTests to true")
	}

	testsFile, err := loadTestCasesFromServer[CaptchaSiteverifyTestCasesFile](CaptchaSiteverifyTestCasesJSONEndpoint)
	if err != nil {
		t.Fatalf("failed to load test cases from mock server: %v", err)
	}

	for _, test := range testsFile.Tests {
		t.Run(test.Name, func(t *testing.T) {
			frcClient, err := NewClient(
				WithAPIKey("YOUR_API_KEY"),
				WithSitekey("YOUR_SITE_KEY"),
				WithAPIEndpoint(MockServerURL),
				WithStrictMode(test.Strict),
			)
			if err != nil {
				t.Fatalf("failed to create Friendly Captcha client: %v", err)
			}
			result := frcClient.VerifyCaptchaResponse(context.TODO(), test.Response)

			shouldAccept := result.ShouldAccept()
			expectedShouldAccept := test.Expectation.ShouldAccept

			t.Logf("Result: %+v", result) // Helps with debugging in case of failure

			assert.Equal(
				t,
				expectedShouldAccept,
				shouldAccept,
				fmt.Sprintf("Expected shouldAccept to be: %v, got: %v", expectedShouldAccept, shouldAccept),
			)

			assert.Equal(
				t,
				test.Expectation.WasAbleToVerify,
				result.WasAbleToVerify(),
				fmt.Sprintf("Expected WasAbleToVerify to be: %v, got: %v", test.Expectation.WasAbleToVerify, result.WasAbleToVerify()),
			)

			assert.Equal(
				t,
				test.Expectation.IsClientError,
				result.IsErrorDueToClientError(),
				fmt.Sprintf("Expected IsClientError to be: %v, got: %v", test.Expectation.IsClientError, result.IsErrorDueToClientError()),
			)

			if result.Success {
				var expectedResponse SuccessCaptchaSiteverifyResponse
				err := json.Unmarshal(test.SiteverifyResponse, &expectedResponse)
				if err != nil {
					t.Fatalf("Failed to unmarshal expected siteverify response: %v", err)
				}

				exp := expectedResponse.Data
				res := result.response.Data

				assert.Equal(
					t,
					exp.EventID,
					res.EventID,
					"Event ID does not match expected value",
				)

				assert.Equal(
					t,
					exp.Challenge,
					res.Challenge,
					"Challenge data does not match expected value",
				)

				assert.Equal(
					t,
					exp.RiskIntelligence,
					res.RiskIntelligence,
					"Risk Intelligence data does not match expected value",
				)

				// Check two specific fields:
				assert.Equal(
					t,
					exp.RiskIntelligence.V.Client.HeaderUserAgent,
					res.RiskIntelligence.V.Client.HeaderUserAgent,
				)
				assert.Equal(
					t,
					exp.RiskIntelligence.V.Client.Browser.V.ID,
					res.RiskIntelligence.V.Client.Browser.V.ID,
				)

				if exp.RiskIntelligence.Valid {
					assert.Contains(
						t,
						string(res.RiskIntelligenceRaw.V),
						"header_user_agent",
					)
				}
			}
		})
	}
}

func TestSDKWithMockServerRiskIntelligenceRetrieve(t *testing.T) {
	if !runTestSDKMockServerTests {
		t.Skip("Skipping SDK tests, run the mock server first and set runTestSDKTests to true")
	}

	testsFile, err := loadTestCasesFromServer[RiskIntelligenceRetrieveTestCasesFile](RiskIntelligenceRetrieveTestCasesEndpoint)
	if err != nil {
		t.Fatalf("failed to load test cases from mock server: %v", err)
	}

	for _, test := range testsFile.Tests {
		t.Run(test.Name, func(t *testing.T) {
			frcClient, err := NewClient(
				WithAPIKey("YOUR_API_KEY"),
				WithAPIEndpoint(MockServerURL),
			)
			if err != nil {
				t.Fatalf("failed to create Friendly Captcha client: %v", err)
			}

			result := frcClient.RetrieveRiskIntelligence(context.TODO(), test.Token)
			t.Logf("Result: %+v", result)

			assert.Equal(
				t,
				test.Expectation.WasAbleToRetrieve,
				result.WasAbleToRetrieve(),
				fmt.Sprintf(
					"Expected WasAbleToRetrieve to be: %v, got: %v",
					test.Expectation.WasAbleToRetrieve,
					result.WasAbleToRetrieve(),
				),
			)

			assert.Equal(
				t,
				test.Expectation.IsClientError,
				result.IsErrorDueToClientError(),
				fmt.Sprintf(
					"Expected IsClientError to be: %v, got: %v",
					test.Expectation.IsClientError,
					result.IsErrorDueToClientError(),
				),
			)

			if result.Success {
				var expectedResponse SuccessRiskIntelligenceRetrieveResponse
				err := json.Unmarshal(test.RetrieveResponse, &expectedResponse)
				if err != nil {
					t.Fatalf("Failed to unmarshal expected retrieve response: %v", err)
				}

				exp := expectedResponse.Data
				res := result.response.Data
				if !assert.NotNil(t, res, "Retrieve response data should not be nil on success") {
					return
				}

				assert.Equal(
					t,
					exp.EventID,
					res.EventID,
					"Event ID does not match expected value",
				)

				assert.Equal(
					t,
					exp.RiskIntelligence,
					res.RiskIntelligence,
					"Risk Intelligence data does not match expected value",
				)

				assert.Equal(t, exp.Token, res.Token, "Retrieve token does not match expected value")

				if exp.RiskIntelligence.Valid {
					assert.Equal(
						t,
						exp.RiskIntelligence.V.Client.HeaderUserAgent,
						res.RiskIntelligence.V.Client.HeaderUserAgent,
					)

					if exp.RiskIntelligence.V.Client.Browser.Valid && res.RiskIntelligence.V.Client.Browser.Valid {
						assert.Equal(
							t,
							exp.RiskIntelligence.V.Client.Browser.V.ID,
							res.RiskIntelligence.V.Client.Browser.V.ID,
						)
					}
				}
			}
		})
	}
}

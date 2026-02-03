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
	MockServerURL         = "http://localhost:1090"
	SiteverifyEndpoint    = "/api/v2/captcha/siteverify"
	TestCasesJSONEndpoint = "/api/v1/tests"
)

var runTestSDKMockServerTests = true

type TestCasesFile struct {
	Version int        `json:"version"`
	Tests   []TestCase `json:"tests"`
}

type TestCase struct {
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

type SuccessSiteverifyResponse struct {
	Data VerifyResponseData `json:"data"`
}

func loadTestCasesFromServer() (TestCasesFile, error) {
	resp, err := http.Get(MockServerURL + TestCasesJSONEndpoint)
	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) {
			return TestCasesFile{}, fmt.Errorf("sdk mock testing server is not running, please run it first")
		}
		return TestCasesFile{}, fmt.Errorf("failed to load test cases from mock server: %w", err)
	}
	defer resp.Body.Close()

	var testCases TestCasesFile
	err = json.NewDecoder(resp.Body).Decode(&testCases)
	if err != nil {
		return TestCasesFile{}, fmt.Errorf("failed to decode test cases from mock server: %w", err)
	}

	return testCases, nil
}

func TestSDKWithMockServer(t *testing.T) {
	if !runTestSDKMockServerTests {
		t.Skip("Skipping SDK tests, run the mock server first and set runTestSDKTests to true")
	}

	// Load your test cases
	testsFile, err := loadTestCasesFromServer()
	if err != nil {
		t.Fatalf("failed to load test cases from mock server: %v", err)
	}

	for _, test := range testsFile.Tests {
		t.Run(test.Name, func(t *testing.T) {
			frcClient, err := NewClient(
				WithAPIKey("YOUR_API_KEY"),
				WithSitekey("YOUR_SITE_KEY"),
				WithSiteverifyEndpoint(MockServerURL+SiteverifyEndpoint),
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
				var expectedResponse SuccessSiteverifyResponse
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

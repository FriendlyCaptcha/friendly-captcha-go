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
		ShouldAccept bool `json:"should_accept"`
	} `json:"expectation"`
	Strict bool `json:"strict"`
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
		})
	}
}

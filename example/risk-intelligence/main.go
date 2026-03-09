package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	friendlycaptcha "github.com/friendlycaptcha/friendly-captcha-go"
)

type templateData struct {
	Message             string
	Sitekey             string
	AgentEndpoint       string
	RiskToken           string
	RequestStatus       int
	TokenTimestamp      string
	TokenExpiresAt      string
	TokenNumUses        int64
	TokenOrigin         string
	RiskIntelligenceRaw string
}

func main() {
	sitekey := os.Getenv("FRC_SITEKEY")
	apikey := os.Getenv("FRC_APIKEY")
	apiEndpoint := os.Getenv("FRC_API_ENDPOINT")
	agentEndpoint := os.Getenv("FRC_AGENT_ENDPOINT")

	if sitekey == "" || apikey == "" {
		log.Fatalf("Please set FRC_SITEKEY and FRC_APIKEY before running this example.")
	}

	opts := []friendlycaptcha.ClientOption{
		friendlycaptcha.WithAPIKey(apikey),
		friendlycaptcha.WithSitekey(sitekey),
	}
	if apiEndpoint != "" {
		opts = append(opts, friendlycaptcha.WithAPIEndpoint(apiEndpoint))
	}

	frcClient, err := friendlycaptcha.NewClient(opts...)
	if err != nil {
		log.Fatalf("Failed to create Friendly Captcha client: %s", err)
	}

	tmpl := template.Must(template.ParseFiles("demo.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			renderTemplate(w, tmpl, templateData{
				Sitekey:       sitekey,
				AgentEndpoint: agentEndpoint,
			})
			return
		}

		riskToken := strings.TrimSpace(r.FormValue("frc-risk-intelligence-token"))
		if riskToken == "" {
			renderTemplate(w, tmpl, templateData{
				Message:       "No risk intelligence token found",
				Sitekey:       sitekey,
				AgentEndpoint: agentEndpoint,
			})
			return
		}

		result := frcClient.RetrieveRiskIntelligence(r.Context(), riskToken)
		data := templateData{
			Sitekey:       sitekey,
			AgentEndpoint: agentEndpoint,
			RiskToken:     riskToken,
			RequestStatus: result.HTTPStatusCode(),
		}

		if !result.WasAbleToRetrieve() {
			data.Message = fmt.Sprintf("Risk intelligence retrieval failed: %v", result.RequestError())
			renderTemplate(w, tmpl, data)
			return
		}

		res := result.Response()

		if !result.IsValid() {
			data.Message = fmt.Sprintf("Risk intelligence token is invalid: %v (%v)", res.Error.Detail, res.Error.ErrorCode)
			renderTemplate(w, tmpl, data)
			return
		}

		if !res.Data.RiskIntelligence.Valid {
			// This should never happen since /retrieve should always return
			// risk intelligence data if the token is valid and request is successful.
			// But it's good practice to handle it just in case.
			data.Message = "Token was valid, but risk intelligence data was not returned."
			data.TokenTimestamp = res.Data.Token.Timestamp.Format(time.RFC3339)
			data.TokenExpiresAt = res.Data.Token.ExpiresAt.Format(time.RFC3339)
			data.TokenNumUses = res.Data.Token.NumUses
			data.TokenOrigin = res.Data.Token.Origin
			renderTemplate(w, tmpl, data)
			return
		}

		prettyJSON, err := json.MarshalIndent(res.Data.RiskIntelligence.V, "", "  ")
		if err != nil {
			data.Message = fmt.Sprintf("Retrieved data, but failed to format JSON: %v", err)
			renderTemplate(w, tmpl, data)
			return
		}

		data.Message = "Retrieved risk intelligence data successfully."
		data.TokenTimestamp = res.Data.Token.Timestamp.Format(time.RFC3339)
		data.TokenExpiresAt = res.Data.Token.ExpiresAt.Format(time.RFC3339)
		data.TokenNumUses = res.Data.Token.NumUses
		data.TokenOrigin = res.Data.Token.Origin
		data.RiskIntelligenceRaw = string(prettyJSON)
		renderTemplate(w, tmpl, data)
	})

	log.Printf("Starting server on http://localhost:8845")
	log.Fatal(http.ListenAndServe(":8845", nil))
}

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, data templateData) {
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

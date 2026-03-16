package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	friendlycaptcha "github.com/friendlycaptcha/friendly-captcha-go"
)

type formMessage struct {
	Subject string
	Message string
}

type templateData struct {
	Message        string
	Sitekey        string
	WidgetEndpoint string
}

func main() {
	sitekey := os.Getenv("FRC_SITEKEY")
	apikey := os.Getenv("FRC_APIKEY")

	// Optionally we can pass in custom endpoints to be used, such as "eu".
	apiEndpoint := os.Getenv("FRC_API_ENDPOINT")
	widgetEndpoint := os.Getenv("FRC_WIDGET_ENDPOINT")

	if sitekey == "" || apikey == "" {
		log.Fatalf("Please set the FRC_SITEKEY and FRC_APIKEY environment values before running this example to your Friendly Captcha sitekey and API key respectively.")
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
				Sitekey:        sitekey,
				WidgetEndpoint: widgetEndpoint,
			})
			return
		}

		form := formMessage{
			Subject: r.FormValue("subject"),
			Message: r.FormValue("message"),
		}

		retrieveRiskIntelligenceIfAvailable(r.Context(), frcClient, r.FormValue("frc-risk-intelligence-token"))

		solution := r.FormValue("frc-captcha-response")
		result := frcClient.VerifyCaptchaResponse(r.Context(), solution)

		if !result.WasAbleToVerify() {
			// In this case we were not actually able to verify the response embedded in the form, but we may still want to accept it.
			// It could mean there is a network issue or that the service is down. In those cases you generally want to accept submissions anyhow.
			// That's why we use `ShouldAccept()` below to actually accept or reject the form submission. It will return true in these cases.

			if result.IsErrorDueToClientError() {
				// Something is wrong with our configuration, check your API key!
				// Send yourself an alert to fix this! Your site is unprotected until you fix this.
				log.Printf("CAPTCHA CONFIG ERROR: %s\n", result.RequestError())
			} else {
				log.Printf("Failed to verify captcha response: %s\n", result.RequestError())
			}
		}

		if !result.ShouldAccept() {
			renderTemplate(w, tmpl, templateData{
				Message:        "❌ Anti-robot check failed, please try again.",
				Sitekey:        sitekey,
				WidgetEndpoint: widgetEndpoint,
			})
			return
		}

		// The captcha was OK, process the form.
		_ = form // Normally we would use the form data here and submit it to our database.

		renderTemplate(w, tmpl, templateData{
			Message:        "✅ Your message has been submitted successfully.",
			Sitekey:        sitekey,
			WidgetEndpoint: widgetEndpoint,
		})
	})

	log.Printf("Starting server on localhost port 8844 (http://localhost:8844)")
	log.Fatal(http.ListenAndServe(":8844", nil))
}

func retrieveRiskIntelligenceIfAvailable(ctx context.Context, frcClient *friendlycaptcha.Client, token string) {
	token = strings.TrimSpace(token)
	if token == "" {
		log.Printf("No risk intelligence token found in form data, skipping retrieval.")
		return
	}

	result := frcClient.RetrieveRiskIntelligence(ctx, token)
	if !result.WasAbleToRetrieve() {
		log.Printf("Failed to retrieve risk intelligence: %v", result.RequestError())
		return
	}

	if !result.IsValid() {
		log.Printf("Risk intelligence token is invalid: %+v", result.Response().Error)
		return
	}

	response := result.Response()
	if !response.Data.RiskIntelligence.Valid {
		log.Printf("Risk intelligence retrieval succeeded, but no risk intelligence data was returned.")
		return
	}

	prettyJSON, err := json.MarshalIndent(response.Data.RiskIntelligence.V, "", "  ")
	if err != nil {
		log.Printf("Retrieved risk intelligence but failed to format JSON: %v", err)
		return
	}

	log.Printf("Risk Intelligence Data:\n%s", prettyJSON)
	log.Printf("Token data: %+v", response.Data.Token)
}

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, data templateData) {
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

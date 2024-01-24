package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

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
	siteverifyEndpoint := os.Getenv("FRC_SITEVERIFY_ENDPOINT")
	widgetEndpoint := os.Getenv("FRC_WIDGET_ENDPOINT")

	if sitekey == "" || apikey == "" {
		log.Fatalf("Please set the FRC_SITEKEY and FRC_APIKEY environment values before running this example to your Friendly Captcha sitekey and API key respectively.")
	}

	opts := []friendlycaptcha.ClientOption{
		friendlycaptcha.WithAPIKey(apikey),
		friendlycaptcha.WithSitekey(sitekey),
	}
	if siteverifyEndpoint != "" {
		opts = append(opts, friendlycaptcha.WithSiteverifyEndpoint(siteverifyEndpoint)) // optional, defaults to "global"
	}
	frcClient, err := friendlycaptcha.NewClient(opts...)
	if err != nil {
		log.Fatalf("Failed to create Friendly Captcha client: %s", err)
	}

	tmpl := template.Must(template.ParseFiles("demo.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// GET - the user is requesting the form, not submitting it.
		if r.Method != http.MethodPost {
			err := tmpl.Execute(w, templateData{
				Message:        "",
				Sitekey:        sitekey,
				WidgetEndpoint: widgetEndpoint,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		form := formMessage{
			Subject: r.FormValue("subject"),
			Message: r.FormValue("message"),
		}

		solution := r.FormValue("frc-captcha-response")
		result := frcClient.VerifyCaptchaResponse(r.Context(), solution)

		if !result.WasAbleToVerify() {
			// In this case we were not actually able to verify the response embedded in the form, but we may still want to accept it.
			// It could mean there is a network issue or that the service is down. In those cases you generally want to accept submissions anyhow.
			// That's why we use `shouldAccept()` below to actually accept or reject the form submission. It will return true in these cases.

			if result.IsErrorDueToClientError() {
				// Something is wrong with our configuration, check your API key!
				// Send yourself an alert to fix this! Your site is unprotected until you fix this.
				log.Printf("CAPTCHA CONFIG ERROR: %s\n", result.RequestError())
			} else {
				log.Printf("Failed to verify captcha response: %s\n", result.RequestError())
			}
		}

		if !result.ShouldAccept() {
			err := tmpl.Execute(w, templateData{
				Message:        "❌ Anti-robot check failed, please try again.",
				Sitekey:        sitekey,
				WidgetEndpoint: widgetEndpoint,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		// The captcha was OK, process the form.
		_ = form // Normally we would use the form data here and submit it to our database.

		err := tmpl.Execute(w, templateData{
			Message:        "✅ Your message has been submitted successfully.",
			Sitekey:        sitekey,
			WidgetEndpoint: widgetEndpoint,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Printf("Starting server on localhost port 8844 (http://localhost:8844)")
	http.ListenAndServe(":8844", nil)
}

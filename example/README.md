# Friendly Captcha Go Example

This application integrates Friendly Captcha for form submissions using Go.
It verifies captcha responses and retrieves risk intelligence (if enabled on the application) from the same form flow.

### Requirements

- Go
- Your Friendly Captcha API key and sitekey.

### Start the application

- Set environment variables and start the application

> NOTE: `FRC_API_ENDPOINT` and `FRC_WIDGET_ENDPOINT` are optional. If not set, default values will be used. You can also use `global` or `eu` as shorthands for both.

```bash
FRC_APIKEY=<your API key> FRC_SITEKEY=<your sitekey> FRC_API_ENDPOINT=<api endpoint> FRC_WIDGET_ENDPOINT=<widget endpoint> go run .
```

## Usage

Navigate to http://localhost:8844/ in your browser.
Fill out the form and submit. The backend verifies the captcha and also retrieves risk intelligence data when a risk intelligence token is available.

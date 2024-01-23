# Friendly Captcha Go Example

This application integrates Friendly Captcha for form submissions using Go.

### Requirements

- Go
- Your Friendly Captcha API key and site key.

### Start the application

- Setup env variables and start the application

> NOTE: `FRC_SITEVERIFY_ENDPOINT` and `FRC_WIDGET_ENDPOINT` are optional. If not set, the default values will be used. You can also use `global` or `eu` as shorthands for both.

```bash
FRC_APIKEY=<your api key> FRC_SITEKEY=<your site key> FRC_SITEVERIFY_ENDPOINT=<siteverify endpoint> FRC_WIDGET_ENDPOINT=<widget endpoint> go run main.go
```

# Usage

Navigate to http://localhost:8844/ in your browser.
Fill out the form and submit. The Friendly Captcha verification will protect the form from bots.

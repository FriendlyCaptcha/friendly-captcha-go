# Friendly Captcha Go Risk Intelligence Example

This example demonstrates server-side risk intelligence retrieval with `RetrieveRiskIntelligence`.

### Requirements

- Go
- Your Friendly Captcha API key and sitekey.

### Start the application

- Setup env variables and start the application

> NOTE: `FRC_API_ENDPOINT` and `FRC_AGENT_ENDPOINT` are optional. If not set, the default values will be used. You can also use `global` or `eu` as shorthands for both.

```bash
FRC_APIKEY=<your API key> FRC_SITEKEY=<your sitekey> FRC_API_ENDPOINT=<api endpoint> FRC_AGENT_ENDPOINT=<agent endpoint> go run main.go

```

## Usage

Navigate to http://localhost:8845/ in your browser.
The token generation starts automatically. Submit the form to retrieve the risk intelligence data server-side.
Tokens are cached in the browser for the duration of their validity period so refreshing the page will not regenerate the token.
You can regenerate the token by clicking the "Regenerate Token" button.

# Friendly Captcha Go SDK

A Go client for the [Friendly Captcha](https://friendlycaptcha.com) service. This client allows for easy integration and verification of captcha responses with the Friendly Captcha API.

> This library is for [Friendly Captcha V2](https://developer.friendlycaptcha.com) only. If you are looking for V1, look [here](https://docs.friendlycaptcha.com)

## Installation

```shell
go get github.com/friendlycaptcha/friendly-captcha-go
```

## Usage

Below are some basic examples of how to use the client.

For a more detailed example, take a look at the [example](./example) directory.

### Initialization

```go
import friendlycaptcha "github.com/friendlycaptcha/friendly-captcha-go"
...
opts := []friendlycaptcha.ClientOption{
		friendlycaptcha.WithAPIKey("YOUR_API_KEY"),
		friendlycaptcha.WithSitekey("YOUR_SITEKEY"),
}
frcClient, err := friendlycaptcha.NewClient(opts...)
if err != nil {
    // handle possible configuration error
}
```

### Verifying a Captcha Response

After calling `VerifyCaptchaResponse` with the captcha response there are two functions on the result object that you should check:

- `WasAbleToVerify()` indicates whether we were able to verify the captcha response. This will be `false` in case there was an issue with the network/our service or if there was a mistake in the configuration.
- `ShouldAccept()` indicates whether the captcha response was correct. If the client is running in non-strict mode (default) and `WasAbleToVerify()` returned `false`, this will be `true`.

Below are some examples of this behaviour.

#### Verifying a correct captcha response without issues when veryfing:

```go
result := frcClient.VerifyCaptchaResponse(context.TODO(), "CORRECT_CAPTCHA_RESPONSE_HERE")
fmt.Println(result.WasAbleToVerify()) // true
fmt.Println(result.ShouldAccept()) // true
```

#### Verifying an incorrect captcha response without issues when veryfing:

```go
result := frcClient.VerifyCaptchaResponse(context.TODO(), "INCORRECT_CAPTCHA_RESPONSE_HERE")
fmt.Println(result.WasAbleToVerify()) // true
fmt.Println(result.ShouldAccept()) // false
```

#### Verifying an incorrect captcha response with issues (network issues or bad configuration) when veryfing in non-strict mode (default):

```go
result := frcClient.VerifyCaptchaResponse(context.TODO(), "INCORRECT_CAPTCHA_RESPONSE_HERE")
fmt.Println(result.WasAbleToVerify()) // false
fmt.Println(result.ShouldAccept()) // true
```

#### Verifying an incorrect captcha response with issues (network/service issues or bad configuration) when veryfing in strict mode:

```go
frcClient, _ := friendlycaptcha.NewClient(
    ...
    friendlycaptcha.WithStrictMode(true),
)
result := frcClient.VerifyCaptchaResponse(context.TODO(), "INCORRECT_CAPTCHA_RESPONSE_HERE")
fmt.Println(result.WasAbleToVerify()) // false
fmt.Println(result.ShouldAccept()) // false
```

### Configuration

The client offers several configuration options:

- **WithAPIKey**: Your Friendly Captcha API key.
- **WithSitekey**: Your Friendly Captcha sitekey.
- **WithStrictMode**: (Optional) In case the client was not able to verify the captcha response at all (for example if there is a network failure or a mistake in configuration), by default the `VerifyCaptchaResponse` returns `True` regardless. By passing `WithStrictMode(true)`, it will return `false` instead: every response needs to be strictly verified.
- **WithSiteverifyEndpoint**: (Optional) The endpoint URL for the site verification API. Shorthands `eu` or `global` are also accepted. Default is `global`.

## Development

### Run the tests

First run the SDK Test server, then run `go test`.

```shell
docker run -p 1090:1090 friendlycaptcha/sdk-testserver:latest

go test -v -tags=sdkintegration ./...
```

## License

Open source under [MIT](./LICENSE).

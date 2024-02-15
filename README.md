# friendly-captcha-go

A Go client for the [Friendly Captcha](https://friendlycaptcha.com) service. This client allows for easy integration and verification of captcha responses with the Friendly Captcha API.

> This library is for [Friendly Captcha v2](https://developer.friendlycaptcha.com) only. If you are looking for V1, look [here](https://docs.friendlycaptcha.com)

## Installation

**Install using [NPM](https://npmjs.com/)**

```shell
go get github.com/friendlycaptcha/friendly-captcha-go
```

## Usage

First configure and create a SDK client

```go
// TODO: add an example here.
```

## Development

### Run the tests
First run the SDK Test server, then run `go test`.
```shell
docker run -p 1090:1090 friendlycaptcha/sdk-testserver:latest

go test ./...
```

## License

Open source under [MIT](./LICENSE).

package schmokin

type HTTPClient interface {
	Execute(request Request) Result
}

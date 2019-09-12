package schmokin

type Reporter interface {
	Execute(result Result)
}

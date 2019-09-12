package schmokin

type Result struct {
	Headers map[string][]string
}

func NewResult() Result {
	return Result{
		Headers: map[string][]string{},
	}
}

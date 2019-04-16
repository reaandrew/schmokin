package main

type SchmokinResult struct {
	Results ResultCollection
	Error   error
}

func (self SchmokinResult) Success() bool {
	return self.Results.Success()
}

type InsufficientArguments SchmokinResult

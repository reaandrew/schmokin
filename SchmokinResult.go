package main

type SchmokinResult struct {
	Results ResultCollection
}

func (self SchmokinResult) Success() bool {
	return self.Results.Success()
}

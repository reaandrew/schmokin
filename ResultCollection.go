package main

type ResultCollection []Result

func (instance ResultCollection) Success() bool {
	for _, result := range instance {
		if !result.Success {
			return false
		}
	}
	return true
}

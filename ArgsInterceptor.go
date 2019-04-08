package main

import (
	"io/ioutil"
	"os"
	"strings"
)

type ArgsInterceptor struct {
	state State
}

func (self ArgsInterceptor) Intercept(args []string) []string {
	current := 0
	for current < len(args) {
		if len(args) < current+2 {
			break
		}
		switch args[current] {
		case "-d":
			value := args[current+1]
			if strings.HasPrefix(value, "@") {
				data, _ := ioutil.ReadFile(value[1:])
				f, _ := os.Create("schmokin.payload")
				defer f.Close()
				f.Write([]byte(self.state.Replace(string(data))))
				args[current+1] = "@schmokin.payload"
			} else {
				args[current+1] = self.state.Replace(value)
			}
			current += 1
		default:
			value := args[current+1]
			args[current+1] = self.state.Replace(value)
			current += 1
			break
		}
	}
	return args
}

func CreateArgsInterceptor(state State) ArgsInterceptor {
	return ArgsInterceptor{
		state: state,
	}
}

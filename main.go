package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	CommitHash string
	Version    string
	BuildTime  string
)

func main() {
	f, _ := os.Open("sample2.yml")
	fmt.Println(CommitHash)
	fmt.Println(Version)
	fmt.Println(BuildTime)
	scanner := bufio.NewScanner(f)
	request := []string{}
	data := []string{}
	line := 0
	setData := false
	for scanner.Scan() {
		lineContent := scanner.Text()
		if line > 0 && lineContent == "---" {
			setData = true
		}

		if !setData {
			request = append(request, lineContent)
		} else {
			if lineContent != "---" {
				data = append(data, lineContent)
			}
		}
		line++
	}

	fmt.Println(strings.Join(request, "\n"))
	fmt.Println(data)
}

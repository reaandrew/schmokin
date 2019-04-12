package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

func WriteFile(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(data)
	return nil
}

func ReadFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ReadLines(file *os.File, visitor func(line string)) {
	reader := bufio.NewReader(file)
	for {
		var buffer bytes.Buffer

		var l []byte
		var err error
		var isPrefix bool
		for {
			l, isPrefix, err = reader.ReadLine()
			buffer.Write(l)

			// If we've reached the end of the line, stop reading.
			if !isPrefix {
				break
			}

			// If we're just at the EOF, break
			if err != nil {
				break
			}
		}

		if err == io.EOF {
			break
		}

		line := buffer.String()
		if line != "" {
			visitor(line)
		}
	}
}

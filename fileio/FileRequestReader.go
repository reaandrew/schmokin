package fileio

import "io/ioutil"

type FileRequestReader struct {
	Path string
}

func (f FileRequestReader) Read() string {
	data, err := ioutil.ReadFile(f.Path)
	if err != nil {
		panic(err)
	}
	return string(data)
}

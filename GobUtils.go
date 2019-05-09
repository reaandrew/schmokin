package main

import (
	"encoding/gob"
	"os"
)

func WriteGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	defer file.Close()
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	return err
}

func ReadGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	defer file.Close()
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	return err
}

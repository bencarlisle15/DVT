package main

import (
	"fmt"
	"io/ioutil"
)

func AddData(data []byte, hash string) {
	filename := fmt.Sprintf("data/%s", hash)
	_ = ioutil.WriteFile(filename, data, 777)
}

func GetData(hash string) []byte {
	filename := fmt.Sprintf("data/%s", hash)
	data, _ := ioutil.ReadFile(filename)
	return data
}
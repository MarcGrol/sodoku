package main

import (
	"encoding/json"
	"io"
)

type Response struct {
	Error     *ErrorDescriptor `json:"error"`
	Solutions []Solution       `json:"solutions"`
}

type ErrorDescriptor struct {
	Message string `json:"message"`
}

type Exercise struct {
	Steps []Step `json:"steps"`
}

type Solution struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

func FromJson(reader io.Reader) (*Exercise, error) {
	exerc := Exercise{}
	err := json.NewDecoder(reader).Decode(&exerc)
	if err != nil {
		return nil, err
	}
	return &exerc, nil
}

func (resp Response) ToJson(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(resp)
}

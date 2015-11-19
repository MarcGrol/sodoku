package main

import (
	"encoding/json"
	"io"
)

type Response struct {
	Error     *ErrorDescriptor `json:"error"`
	Solutions []Game           `json:"solutions"`
}

type ErrorDescriptor struct {
	Message string `json:"message"`
}

type Game struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	X       int  `json:"x"`
	Y       int  `json:"y"`
	Z       int  `json:"z"`
	Initial bool `json:"initial"`
	IsGuess bool `json:"isGuess"`
}

func FromJson(reader io.Reader) (*Game, error) {
	game := Game{}
	err := json.NewDecoder(reader).Decode(&game)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (resp Response) ToJson(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(resp)
}

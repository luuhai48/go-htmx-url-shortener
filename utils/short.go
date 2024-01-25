package utils

import (
	"log"

	"github.com/teris-io/shortid"
)

var Generator *shortid.Shortid

func init() {
	generator, err := shortid.New(1, shortid.DefaultABC, 2343)
	if err != nil {
		log.Fatal(err)
	}
	Generator = generator
}

func GenShortID() (string, error) {
	id, err := Generator.Generate()
	if err != nil {
		return "", err
	}
	return id, nil
}

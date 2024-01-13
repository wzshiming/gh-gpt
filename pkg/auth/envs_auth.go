package auth

import (
	"os"
)

type envsAuth struct{}

func Envs() Auth {
	return envsAuth{}
}

func (envsAuth) GetToken() (string, error) {
	token, ok := os.LookupEnv("GH_TOKEN")
	if !ok {
		return "", ErrNotExists
	}
	return token, nil
}

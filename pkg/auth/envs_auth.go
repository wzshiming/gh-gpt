package auth

import (
	"context"
	"os"
)

type envsAuth struct{}

func Envs() Auth {
	return envsAuth{}
}

func (envsAuth) GetToken(ctx context.Context) (string, error) {
	token, ok := os.LookupEnv("GH_COPILOT_TOKEN")
	if !ok {
		token, ok = os.LookupEnv("GH_TOKEN")
		if !ok {
			return "", ErrNotExists
		}
	}
	return token, nil
}

package auth

import (
	"encoding/json"
	"os"
	"path"
	"runtime"
)

type hostsAuth struct{}

func Hosts() Auth {
	return hostsAuth{}
}

func (hostsAuth) GetToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := ".config"
	if runtime.GOOS == "windows" {
		configPath = "AppData/Local"
	}
	p := path.Join(home, configPath, "github-copilot/hosts.json")
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrNotExists
		}
		return "", err
	}

	var h hosts
	err = json.Unmarshal(data, &h)
	if err != nil {
		return "", err
	}
	return h.Github.OauthToken, nil
}

type hosts struct {
	Github *hostAuth `json:"github.com"`
}

type hostAuth struct {
	User       string `json:"user"`
	OauthToken string `json:"oauth_token"`
}

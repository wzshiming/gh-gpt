package auth

import (
	"os"
)

var ErrNotExists = os.ErrNotExist

type Auth interface {
	GetToken() (string, error)
}

package util

import (
	"os/user"

	"harkd/errors"
)

// GetUserHomeDir provide the home directory of the current system user.
func GetUserHomeDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", errors.ErrUserLookup{err}
	}

	return currentUser.HomeDir, nil
}

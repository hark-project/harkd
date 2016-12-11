package context

import (
	"path/filepath"

	"harkd/dal"
	"harkd/util"
)

const dalFileName = "hark.json"

// Factory is an interface that can provide a hark Context.
type Factory interface {
	GetContext() Context
}

// HomeDirFactory returns a Factory providing a Context
// based on the home directory of the current user.
//
// This is the most simple form of storage available.
func HomeDirFactory() (Factory, error) {
	homeDir, err := util.GetUserHomeDir()
	if err != nil {
		return nil, err
	}

	// Construct a dal
	//
	// The dal is kept in the state of the factory - rather than instantiated on
	// demand when creating a Context - so that its locking facilities can be
	// shared across the app.
	d, err := dal.NewJSONFileDal(dalFilePath(homeDir))
	if err != nil {
		return nil, err
	}

	return dirFactory{homeDir, d}, nil
}

type dirFactory struct {
	dir string
	dal dal.Dal
}

func dalFilePath(contextDir string) string {
	return filepath.Join(contextDir, dalFileName)
}

func (d dirFactory) GetContext() Context {
	return dirContext{d.dir, d.dal}
}

package dal

import (
	"encoding/json"
	"fmt"
	"os"

	"harkd/core"
	"harkd/errors"
	"harkd/util"
	"harkd/util/fs"
)

const jsonFileDalFileMode = 0644
const lockFileSuffix = ".lock"

// NewJSONFileDal returns a DAL which is backed by state in a simple flat JSON
// file.
func NewJSONFileDal(filename string) (Dal, error) {
	lockFileName := filename + lockFileSuffix
	lock, err := util.NewLock(lockFileName)
	if err != nil {
		return nil, err
	}

	return &jsonFileDal{
		filename,
		fs.NewFilesystem(),
		lock,
	}, nil
}

type jsonFileDal struct {
	filename   string
	fileSystem fs.Filesystem
	util.Lock
}

type jsonFileState struct {
	Machines []core.Machine `json:"machines"`
}

func initializeJSONFileState(fileSys fs.Filesystem, filename string) error {
	_, err := fileSys.Stat(filename)
	if err == nil {
		return nil
	} else if err != nil && !os.IsNotExist(err) {
		// TODO(cera) - wrap this error
		return err
	}

	// Set up our empty state
	var jfs jsonFileState

	// Persist it
	return saveJSONFileState(jfs, fileSys, filename)
}

func loadJSONFileState(fileSys fs.Filesystem, filename string) (jsonFileState, error) {
	var jfs jsonFileState

	// Make sure the state has been initialized
	err := initializeJSONFileState(fileSys, filename)
	if err != nil {
		return jfs, err
	}

	f, err := fileSys.Open(filename)
	if err != nil {
		// TODO(cera) - wrap this error
		return jfs, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	return jfs, dec.Decode(&jfs)
}

func saveJSONFileState(jfs jsonFileState, fileSys fs.Filesystem, filename string) error {
	b, err := json.Marshal(jfs)
	if err != nil {
		return errors.ErrSerialization("serializing state", err)
	}

	if err := fileSys.WriteFile(filename, b, 0644); err != nil {
		return errors.ErrStatePersist(err)
	}

	return nil
}

func (jfd jsonFileDal) loadState() (jsonFileState, error) {
	return loadJSONFileState(jfd.fileSystem, jfd.filename)
}

func (jfd jsonFileDal) withState(fn func(s jsonFileState) error) error {
	s, err := jfd.loadState()
	if err != nil {
		return err
	}

	return fn(s)
}

func (jfd jsonFileDal) withStateLock(fn func(s jsonFileState) error) error {
	return jfd.withLock(func() error {
		return jfd.withState(fn)
	})
}

func (jfd jsonFileDal) withLock(fn func() error) error {
	err := jfd.Lock.Lock()
	if err != nil {
		return err
	}
	defer jfd.Lock.Unlock()

	return fn()
}

func (jfd jsonFileDal) GetMachines() (machines []core.Machine, err error) {
	err = jfd.withState(func(s jsonFileState) error {
		machines = s.Machines
		return nil
	})
	return machines, err
}

func (jfd jsonFileDal) GetMachineByID(machineID string) (core.Machine, error) {
	machines, err := jfd.GetMachines()
	if err != nil {
		return core.Machine{}, err
	}
	for _, m := range machines {
		if m.ID == machineID {
			return m, nil
		}
	}
	return core.Machine{}, errors.ErrMachineNotFound(machineID)
}

func (jfd jsonFileDal) SaveMachine(machine core.Machine) error {
	// get a file lock so that we do not race with other processes or goroutines
	return jfd.withStateLock(func(s jsonFileState) error {
		// First, make sure this ID does not exist already
		for _, m := range s.Machines {
			if m.ID == machine.ID {
				return errors.ErrEntityConflict(fmt.Sprintf("already have machine with id %q", m.ID))
			}
		}

		// Add this machine to the state
		s.Machines = append(s.Machines, machine)

		// Save the state back to the file
		return saveJSONFileState(s, jfd.fileSystem, jfd.filename)
	})
}

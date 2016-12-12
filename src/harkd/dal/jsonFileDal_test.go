package dal

import (
	"bytes"
	"testing"

	"harkd/core"
	"harkd/test/fixtures"

	"github.com/stretchr/testify/require"
)

const fakeDalPath = "/tmp/hark"

func getMockDal(t *testing.T) (Dal, *fixtures.FsFixture) {
	// Construct a DAL
	dal, err := NewJSONFileDal(fakeDalPath)
	require.NoError(t, err)

	// Inject a mock filesystem
	fs := fixtures.NewFsFixture()
	dal.(*jsonFileDal).fileSystem = fs

	return dal, fs
}

var getMachinesTests = []struct {
	name           string
	stateJson      string
	expectMachines []core.Machine
	valid          bool
}{
	{"empty state", "{}", nil, true},
	{"empty machine list", `{"machines":[]}`, []core.Machine{}, true},
	{"invalid type for machine list: 1", `{"machines":{}}`, nil, false},
	{"invalid type for machine list: 2", `{"machines":1`, nil, false},
	{"invalid type for machine list: 3", `{"machines":"a"`, nil, false},
}

func TestJSONFileDalGetMachines(t *testing.T) {
	for _, c := range getMachinesTests {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// Prepare
			dal, fs := getMockDal(t)
			buf := bytes.NewBuffer([]byte(c.stateJson))
			fs.MockOpen.WillReturn = fixtures.NewNopCloser(buf)

			// Execute
			machines, err := dal.GetMachines()

			// Assert
			if c.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			require.Equal(t, c.expectMachines, machines)
		})
	}
}

var getMachineByIDTests = []struct {
	name          string
	stateJson     string
	machineId     string
	expectMachine core.Machine
	valid         bool
}{
	{"empty state", "{}", "foo", core.Machine{}, false},
	{"empty machine list", `{"machines":[]}`, "bar", core.Machine{}, false},
	{"invalid existing state", `abcd`, "foo", core.Machine{}, false},
	{"machines exist but no matching ID", `{"machines":[{"id":"bang"}]}`, "foo", core.Machine{}, false},
	{"have exactly 1 matching machine", `{"machines":[{"id":"foo"}]}`, "foo", core.Machine{ID: "foo"}, true},
	{"have multiple machines, 1 matches", `{"machines":[{"id": "bang"},{"id":"bar"}]}`, "bar", core.Machine{ID: "bar"}, true},
}

func TestJSONFileDalGetMachinebyID(t *testing.T) {
	for _, c := range getMachineByIDTests {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// Prepare
			dal, fs := getMockDal(t)
			buf := bytes.NewBuffer([]byte(c.stateJson))
			fs.MockOpen.WillReturn = fixtures.NewNopCloser(buf)

			// Execute
			machine, err := dal.GetMachineByID(c.machineId)

			// Assert
			if c.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			require.Equal(t, c.expectMachine, machine)
		})
	}
}

var saveMachineTests = []struct {
	name        string
	stateBefore string
	machine     core.Machine
	stateAfter  string
	valid       bool
}{
	{"adding the first machine", "{}", core.Machine{ID: "foo"}, `{"machines":[{"id":"foo","name":"","memoryMB":0}]}`, true},
	{"adding machine with invalid existing state: 1", "[]", core.Machine{ID: "foo"}, "", false},
	{"adding machine with invalid existing state: 2", "abcd", core.Machine{ID: "foo"}, "", false},
	{"duplicate machine", `{"machines":[{"id":"foo"}]}`, core.Machine{ID: "foo"}, "", false},
}

func TestJSONFileDalSaveMachine(t *testing.T) {
	for _, c := range saveMachineTests {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// Prepare
			dal, fs := getMockDal(t)
			buf := bytes.NewBuffer([]byte(c.stateBefore))
			fs.MockOpen.WillReturn = fixtures.NewNopCloser(buf)

			// Execute
			err := dal.SaveMachine(c.machine)

			// Assert
			if c.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			require.Equal(t, c.stateAfter, string(fs.MockWriteFile.CalledWithData))
		})
	}
}

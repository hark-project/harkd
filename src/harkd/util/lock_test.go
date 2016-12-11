package util

import (
	"io/ioutil"
	"os"
	//"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func tempFile(t *testing.T) string {
	tf, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	return tf.Name()
}

func TestLockSimple(t *testing.T) {
	tf := tempFile(t)
	defer os.Remove(tf)

	l, err := NewLock(tf)
	if err != nil {
		t.Fatal(err)
	}

	err = l.Lock()
	require.NoError(t, err)
	err = l.Unlock()
	require.NoError(t, err)
}

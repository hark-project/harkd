package fixtures

import "io"

// NopCloser takes a Reader and converts it into a ReadCloser, with the Close()
// set to be a no-op.
func NewNopCloser(r io.Reader) io.ReadCloser {
	return &NopCloser{r, 0}
}

type NopCloser struct {
	io.Reader

	CloseCalls int
}

func (nc *NopCloser) Close() error {
	nc.CloseCalls++
	return nil
}

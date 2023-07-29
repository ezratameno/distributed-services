package log

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	write = []byte("hello world")
	width = uint64(len(write)) + lenWidth
)

// In this test, we create a store with a temporary file and call two test helpers
// to test appending and reading from the store. Then we create the store again
// and test reading from it again to verify that our service will recover its state
// after a restart.
func TestStoreAppendRead(t *testing.T) {

	f, err := ioutil.TempFile("", "store_append_read_test")
	assert.Nil(t, err)
	defer os.Remove(f.Name())

	s, err := newStore(f)
	assert.Nil(t, err)

	testAppend(t, s)
	testRead(t, s)
	testReadAt(t, s)

	s, err = newStore(f)
	assert.Nil(t, err)

	testRead(t, s)

}

func testAppend(t *testing.T, s *store) {
	t.Helper()
	for i := uint64(1); i < 4; i++ {
		n, pos, err := s.Append(write)
		assert.Nil(t, err)
		assert.Equal(t, pos+n, width*i)
	}
}

func testRead(t *testing.T, s *store) {
	t.Helper()
	var pos uint64
	for i := uint64(1); i < 4; i++ {
		read, err := s.Read(pos)
		assert.Nil(t, err)
		assert.Equal(t, write, read)
		pos += width
	}
}

func testReadAt(t *testing.T, s *store) {
	t.Helper()
	for i, offset := uint64(1), int64(0); i < 4; i++ {
		b := make([]byte, lenWidth)
		n, err := s.ReadAt(b, offset)
		assert.Nil(t, err)
		assert.Equal(t, lenWidth, n)
		offset += int64(n)

		size := enc.Uint64(b)
		b = make([]byte, size)
		n, err = s.ReadAt(b, offset)
		assert.Nil(t, err)
		assert.Equal(t, write, b)
		assert.Equal(t, int(size), n)
		offset += int64(n)
	}
}

func TestStoreClose(t *testing.T) {
	f, err := ioutil.TempFile("", "store_close_test")
	assert.Nil(t, err)
	defer os.ReadFile(f.Name())

	s, err := newStore(f)
	assert.Nil(t, err)

	_, _, err = s.Append(write)
	assert.Nil(t, err)

	f, beforeSize, err := openFile(f.Name())
	assert.Nil(t, err)

	err = s.Close()
	assert.Nil(t, err)

	_, afterSize, err := openFile(f.Name())
	assert.Nil(t, err)

	assert.True(t, afterSize > beforeSize)
}

func openFile(name string) (*os.File, int64, error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, 0, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, 0, err
	}

	return f, fi.Size(), nil
}

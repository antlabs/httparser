package httparser

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ZSplit(t *testing.T) {
	got := 0
	need := 2
	err := Split([]byte("hello,world"), []byte(","), func(v []byte) error {
		switch {
		case bytes.Equal(v, []byte("hello")):
			got++
		case bytes.Equal(v, []byte("world")):
			got++
		}
		return nil
	})
	assert.NoError(t, err)

	assert.Equal(t, need, got)
}

func Test_ZSplit_Error(t *testing.T) {
	err := Split([]byte("hello,world"), []byte(","), func(v []byte) error {
		return errors.New("fail")
	})
	assert.Error(t, err)
}

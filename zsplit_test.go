package httparser

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ZSplit(t *testing.T) {
	got := 0
	need := 2
	Split([]byte("hello,world"), []byte(","), func(v []byte) error {
		switch {
		case bytes.Equal(v, []byte("hello")):
			got++
		case bytes.Equal(v, []byte("world")):
			got++
		}
		return nil
	})

	assert.Equal(t, need, got)
}

func Test_ZSplit_Error(t *testing.T) {
	err := Split([]byte("hello,world"), []byte(","), func(v []byte) error {
		return errors.New("fail")
	})
	assert.Error(t, err)
}

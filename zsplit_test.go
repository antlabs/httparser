package httparser

import (
	"bytes"
	"errors"
	"testing"
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
	if err != nil {
		t.Error(err)
	}

	if got != need {
		t.Errorf("got %d, need %d", got, need)
	}
}

func Test_ZSplit_Error(t *testing.T) {
	err := Split([]byte("hello,world"), []byte(","), func(v []byte) error {
		return errors.New("fail")
	})

	if err == nil {
		t.Error("err should not be nil")
	}
}

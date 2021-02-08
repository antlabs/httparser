package httparser

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Twobuf(t *testing.T) {
	blockSize := 4
	tb := NewTwoBuf(blockSize)

	testData := []string{
		"123456789abcdefg",
		"aaaaaaabbbbbbbbb",
	}

	offset := 1
	for _, v := range testData {
		r := strings.NewReader(v)

		for i := 0; ; i++ {
			buf := tb.Right()
			n, err := r.Read(buf)
			if err == io.EOF {
				break
			}

			if i >= 2 {
				start := i * blockSize

				assert.Equal(t, string(tb.All(n)), v[start-offset:min(int32(start+blockSize), int32(len(v)))])
			} else {
				start := i * blockSize
				assert.Equal(t, string(tb.All(n)), v[start:min(int32(start+blockSize), int32(len(v)))])
			}

			if i != 0 {
				tb.MoveLeft(buf[len(buf)-offset:])
			}
		}

		tb.Reset()
	}
}

func Test_TwobufPanic(t *testing.T) {
	assert.Panics(t, func() {
		tb := NewTwoBuf(4)
		tb.MoveLeft([]byte("12345"))
	})
}

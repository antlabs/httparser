package httparser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Unhex(t *testing.T) {
	for i := 0; i < 256; i++ {
		v := unhex[i]
		switch {
		case i >= '0' && i <= '9':
			assert.Equal(t, v, int8(i-'0'))
		case i >= 'a' && i <= 'f':
			assert.Equal(t, v, int8(i-'a'+10))
		case i >= 'A' && i <= 'F':
			assert.Equal(t, v, int8(i-'A'+10))
		default:
			assert.Equal(t, v, int8(-1), fmt.Sprintf("fail:%c", i))
		}
	}
}

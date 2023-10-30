package httparser

import (
	"testing"
)

func Test_Unhex(t *testing.T) {
	for i := 0; i < 256; i++ {
		v := unhex[i]
		switch {
		case i >= '0' && i <= '9':
			if i-'0' != int(v) {
				t.Fatalf("fail:%c", i)
			}
		case i >= 'a' && i <= 'f':
			if i-'a'+10 != int(v) {
				t.Fatalf("fail:%c", i)
			}

		case i >= 'A' && i <= 'F':
			if i-'A'+10 != int(v) {
				t.Fatalf("fail:%c", i)
			}
		default:
			if v != -1 {
				t.Fatalf("fail:%c", i)
			}

		}
	}
}

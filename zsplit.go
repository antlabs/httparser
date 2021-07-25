package httparser

import "bytes"

func genSplit(s, sep []byte, sepSave int, cb func([]byte) error) error {
	if len(sep) == 0 {
		return nil
	}
	n := bytes.Count(s, sep) + 1

	n--
	i := 0
	for i < n {
		m := bytes.Index(s, sep)
		if m < 0 {
			break
		}
		err := cb(s[: m+sepSave : m+sepSave])
		if err != nil {
			return err
		}
		//a[i] = s[: m+sepSave : m+sepSave]
		s = s[m+len(sep):]
		i++
	}
	return cb(s)
}

// Split 分割字符串
func Split(s, sep []byte, cb func([]byte) error) error { return genSplit(s, sep, 0, cb) }

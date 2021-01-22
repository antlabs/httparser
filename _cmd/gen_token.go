package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
)

//  https://tools.ietf.org/html/rfc7230#section-3.2.6
//  Most HTTP header field values are defined using common syntax
//  components (token, quoted-string, and comment) separated by
//  whitespace or specific delimiting characters.  Delimiters are chosen
//  from the set of US-ASCII visual characters not allowed in a token
//  (DQUOTE and "(),/:;<=>?@[\]{}").
//
//    token          = 1*tchar
//
//    tchar          = "!" / "#" / "$" / "%" / "&" / "'" / "*"
//                   / "+" / "-" / "." / "^" / "_" / "`" / "|" / "~"
//                   / DIGIT / ALPHA
//                   ; any VCHAR, except delimiters

func main() {

	head := `
	package httparser

    // Automatically generated, do not modify
	var token = [256]byte{
	`

	var code bytes.Buffer

	code.WriteString(head)

	for i := 0; i < 256; i++ {
		val := "0"
		switch {
		case i >= '0' && i <= '9':
			fallthrough
		case i >= 'A' && i <= 'Z':
			fallthrough
		case i >= 'a' && i <= 'z':
			val = fmt.Sprintf("%q", i)
		default:
			switch i {
			case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
				val = fmt.Sprintf("%q", i)
			}
		}

		code.WriteString(val)
		if i != 255 {
			code.WriteString(",")
		}

		if i > 0 && i%16 == 0 {
			code.WriteString("\n")
		}
	}

	code.WriteString("}")

	b, err := format.Source(code.Bytes())
	if err != nil {
		panic(err.Error())
	}

	os.Stdout.Write(b)

}

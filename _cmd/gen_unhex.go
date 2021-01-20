package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
)

func main() {
	head := `
	package httparser

    // Automatically generated, do not modify
	var unhex = [256]int8{
	`

	var code bytes.Buffer

	code.WriteString(head)

	for i := 0; i < 256; i++ {
		val := "-1"
		switch {
		case i >= '0' && i <= '9':
			val = string(i)
		case i >= 'a' && i <= 'f':
			val = fmt.Sprintf("%d", i-'a'+10)
		case i >= 'A' && i <= 'F':
			val = fmt.Sprintf("%d", i-'A'+10)
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

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
)

func main() {
	methodsName := []string{
		"GET",
		"HEAD",
		"POST",
		"PUT",
		"DELETE",
		"CONNECT",
		"OPTIONS",
		"TRACE",
		"ACL",
		"BIND",
		"COPY",
		"CHECKOUT",
		"LOCK",
		"UNLOCK",
		"LINK",
		"MKCOL",
		"MOVE",
		"MKACTIVITY",
		"MERGE",
		"M-SEARCH",
		"MKCALENDAR",
		"NOTIFY",
		"PROPFIND",
		"PROPPATCH",
		"PATCH",
		"PURGE",
		"REPORT",
		"REBIND",
		"SUBSCRIBE",
		"SEARCH",
		"SOURCE",
		"UNSUBSCRIBE",
		"UNBIND",
		"UNLINK"}
	methodsType := []string{
		"GET",
		"HEAD",
		"POST",
		"PUT",
		"DELETE",
		"CONNECT",
		"OPTIONS",
		"TRACE",
		"ACL",
		"BIND",
		"COPY",
		"CHECKOUT",
		"LOCK",
		"UNLOCK",
		"LINK",
		"MKCOL",
		"MOVE",
		"MKACTIVITY",
		"MERGE",
		"MSEARCH",
		"MKCALENDAR",
		"NOTIFY",
		"PROPFIND",
		"PROPPATCH",
		"PATCH",
		"PURGE",
		"REPORT",
		"REBIND",
		"SUBSCRIBE",
		"SEARCH",
		"SOURCE",
		"UNSUBSCRIBE",
		"UNBIND",
		"UNLINK"}

	var w io.Writer
	var code bytes.Buffer
	w = &code

	fmt.Fprint(w, `
	switch {
	`)
	for i := range methodsType {
		fmt.Fprintf(w, `case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "%s"):
	p.Method = %s
	`, methodsName[i], methodsType[i])
	}

	fmt.Fprint(w, "}")

	b, err := format.Source(code.Bytes())
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(b)
}

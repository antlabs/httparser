package httparser

type Method int8

const (
	GET Method = iota + 1
	HEAD
	POST
	PUT
	DELETE
	CONNECT
	OPTIONS
	TRACE
	ACL
	BIND
	COPY
	CHECKOUT
	LOCK
	UNLOCK
	LINK
	MKCOL
	MOVE
	MKACTIVITY
	MERGE
	M_SEARCH //M-SEARCH
	MKCALENDAR
	NOTIFY
	PROPFIND
	PROPPATCH
	PATCH
	PURGE
	REPORT
	REBIND
	SUBSCRIBE
	SEARCH
	SOURCE
	UNSUBSCRIBE
	UNBIND
	UNLINK
)
package httparser

// Method 类型 表示http 方法
type Method int8

const (
	// GET 表示GET方法
	GET Method = iota + 1
	// HEAD 表示HEAD方法
	HEAD
	// POST 表示POST方法
	POST
	// PUT 表示PUT方法
	PUT
	// DELETE 表示DELETE方法
	DELETE
	// CONNECT 表示CONNECT方法
	CONNECT
	// OPTIONS 表示OPTIONS方法
	OPTIONS
	// TRACE 表示TRACE方法
	TRACE
	// ACL 表示ACL方法
	ACL
	// BIND 表示BIND方法
	BIND
	// COPY 表示COPY方法
	COPY
	// CHECKOUT 表示CHECKOUT方法
	CHECKOUT
	// LOCK 表示LOCK方法
	LOCK
	// UNLOCK 表示UNLOCK方法
	UNLOCK
	// LINK 表示LINK方法
	LINK
	// MKCOL 表示MKCOL方法
	MKCOL
	// MOVE 表示MOVE方法
	MOVE
	// MKACTIVITY 表示MKACTIVITY方法
	MKACTIVITY
	// MERGE 表示MERGE方法
	MERGE
	// MSEARCH 表示M-SEARCH方法
	MSEARCH
	// MKCALENDAR 表示MKCALENDAR方法
	MKCALENDAR
	// NOTIFY 表示NOTIFY方法
	NOTIFY
	// PROPFIND 表示PROPFIND方法
	PROPFIND
	// PROPPATCH 表示PROPPATCH方法
	PROPPATCH
	// PATCH 表示PATCH方法
	PATCH
	// PURGE 表示PURGE方法
	PURGE
	// REPORT 表示REPORT方法
	REPORT
	// REBIND 表示REBIND方法
	REBIND
	// SUBSCRIBE 表示SUBSCRIBE方法
	SUBSCRIBE
	// SEARCH 表示SEARCH方法
	SEARCH
	// SOURCE 表示SOURCE方法
	SOURCE
	// UNSUBSCRIBE 表示UNSUBSCRIBE方法
	UNSUBSCRIBE
	// UNBIND 表示UNBIND方法
	UNBIND
	// UNLINK 表示UNLINK方法
	UNLINK
)

func (m Method) String() string {
	switch m {
	case GET:
		return "GET"
	case HEAD:
		return "HEAD"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case CONNECT:
		return "CONNECT"
	case OPTIONS:
		return "OPTIONS"
	case TRACE:
		return "TRACE"
	case ACL:
		return "ACL"
	case BIND:
		return "BIND"
	case COPY:
		return "COPY"
	case CHECKOUT:
		return "CHECKOUT"
	case LOCK:
		return "LOCK"
	case UNLOCK:
		return "UNLOCK"
	case LINK:
		return "LINK"
	case MKCOL:
		return "MKCOL"
	case MOVE:
		return "MOVE"
	case MKACTIVITY:
		return "MKACTIVITY"
	case MERGE:
		return "MERGE"
	case MSEARCH:
		return "M-SEARCH"
	case MKCALENDAR:
		return "MKCALENDAR"
	case NOTIFY:
		return "NOTIFY"
	case PROPFIND:
		return "PROPFIND"
	case PROPPATCH:
		return "PROPPATCH"
	case PATCH:
		return "PATCH"
	case PURGE:
		return "PURGE"
	case REPORT:
		return "REPORT"
	case REBIND:
		return "REBIND"
	case SUBSCRIBE:
		return "SUBSCRIBE"
	case SEARCH:
		return "SEARCH"
	case SOURCE:
		return "SOURCE"
	case UNSUBSCRIBE:
		return "UNSUBSCRIBE"
	case UNBIND:
		return "UNBIND"
	case UNLINK:
		return "UNLINK"
	default:
		return "UNKNOWN"
	}
}

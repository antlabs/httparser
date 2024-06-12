// Copyright 2021 guonaihong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httparser

import (
	"reflect"
	"testing"
)

type message struct {
	name           string
	hType          ReqOrRsp
	method         Method
	raw            string
	statusCode     int
	responseStatus string
	//requestPath    string
	requestURL string
	body       string
	//host           string
	//userinfo string
	//port     uint16
	//enum { NONE=0, FIELD, VALUE } last_header_element; //TODO最近移值
	headers         [][2]string
	shouldKeepAlive bool

	//numChunks         int
	//numChunksComplete int
	//chunkLengths      []int

	upgrade string // upgraded body

	httpMajor     uint16
	httpMinor     uint16
	contentLength int32

	messageBeginCbCalled    bool
	headersCompleteCbCalled bool
	messageCompleteCbCalled bool
	//statusCbCalled          bool
	messageCompleteOnEOF bool
	//bodyIsFinal             int
	//allowChunkedLength bool
}

func (m *message) eq(t *testing.T, m2 *message) bool {
	if m.messageCompleteCbCalled != m2.messageCompleteCbCalled {
		t.Errorf("messageCompleteCbCalled: %v != %v", m.messageCompleteCbCalled, m2.messageCompleteCbCalled)
		return false
	}

	if m.method != m2.method {
		t.Errorf("method: %v != %v", m.method, m2.method)
		return false
	}

	if !reflect.DeepEqual(m.headers, m2.headers) {
		t.Errorf("headers: %v\n", m.headers)
		t.Errorf("headers: %v\n", m2.headers)
		return false
	}

	if m.httpMajor != m2.httpMajor {
		t.Errorf("httpMajor: %v != %v", m.httpMajor, m2.httpMajor)
		return false
	}
	if m.httpMinor != m2.httpMinor {
		t.Errorf("httpMinor: %v != %v", m.httpMinor, m2.httpMinor)
		return false
	}

	if m.hType != m2.hType {
		t.Errorf("hType: %v != %v", m.hType, m2.hType)
		return false
	}

	if m.requestURL != m2.requestURL {
		t.Errorf("requestURL: %v != %v", m.requestURL, m2.requestURL)
		return false
	}

	if m.body != m2.body {
		t.Errorf("body: %v != %v", m.body, m2.body)
		return false
	}

	if m.responseStatus != m2.responseStatus {
		t.Errorf("responseStatus: %v != %v", m.responseStatus, m2.responseStatus)
		return false
	}

	if m.statusCode != m2.statusCode {
		t.Errorf("statusCode: %v != %v", m.statusCode, m2.statusCode)
		return false

	}

	if m.upgrade != m2.upgrade {
		t.Errorf("upgrade: %v != %v", m.upgrade, m2.upgrade)
		return false
	}
	return true
}

var requests = []message{
	{
		name:                    "issue 7 1",
		hType:                   REQUEST,
		raw:                     "POST /echo HTTP/1.1\r\nHost: localhost:8080\r\nConnection: close \r\nAccept-Encoding : gzip \r\n\r\n",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		method:                  POST,
		requestURL:              "/echo",
		contentLength:           unused,
		headers: [][2]string{
			{"Host", "localhost:8080"},
			{"Connection", "close "},
			{"Accept-Encoding ", "gzip "},
		},
	},
	{
		name:                    "issue 7 2",
		hType:                   REQUEST,
		raw:                     "POST /echo HTTP/1.1\r\nHost: localhost:8080\r\nConnection: close \r\nContent-Length :  0\r\nAccept-Encoding : gzip \r\n\r\n",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		method:                  POST,
		requestURL:              "/echo",
		contentLength:           unused,
		headers: [][2]string{
			{"Host", "localhost:8080"},
			{"Connection", "close "},
			{"Content-Length ", " 0"},
			{"Accept-Encoding ", "gzip "},
		},
	},
	{
		name:                    "issue 7 3",
		hType:                   REQUEST,
		raw:                     "POST /echo HTTP/1.1\r\nHost: localhost:8080\r\nConnection: close \r\nContent-Length :  5\r\nAccept-Encoding : gzip \r\n\r\nhello",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		method:                  POST,
		requestURL:              "/echo",
		contentLength:           unused,
		body:                    "hello",
		headers: [][2]string{
			{"Host", "localhost:8080"},
			{"Connection", "close "},
			{"Content-Length ", " 5"},
			{"Accept-Encoding ", "gzip "},
		},
	},
	{
		name:                    "issue 7 4",
		hType:                   REQUEST,
		raw:                     "POST / HTTP/1.1\r\nHost: localhost:1235\r\nUser-Agent: Go-http-client/1.1\r\nTransfer-Encoding: chunked\r\nAccept-Encoding: gzip\r\n\r\n4\r\nbody\r\n0\r\n\r\n",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		method:                  POST,
		requestURL:              "/",
		contentLength:           unused,
		body:                    "body",
		headers: [][2]string{
			{"Host", "localhost:1235"},
			{"User-Agent", "Go-http-client/1.1"},
			{"Transfer-Encoding", "chunked"},
			{"Accept-Encoding", "gzip"},
		},
	},
	{
		name:                    "issue 7 5",
		hType:                   REQUEST,
		raw:                     "POST / HTTP/1.1\r\nHost: localhost:1235\r\nUser-Agent: Go-http-client/1.1\r\nTransfer-Encoding: chunked\r\nTrailer: Md5,Size\r\nAccept-Encoding: gzip\r\n\r\n4\r\nbody\r\n0\r\nMd5: 841a2d689ad86bd1611447453c22c6fc\r\nSize: 4\r\n\r\n",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		method:                  POST,
		requestURL:              "/",
		contentLength:           unused,
		body:                    "body",
		headers: [][2]string{
			{"Host", "localhost:1235"},
			{"User-Agent", "Go-http-client/1.1"},
			{"Transfer-Encoding", "chunked"},
			{"Trailer", "Md5,Size"},
			{"Accept-Encoding", "gzip"},
			{"Md5", "841a2d689ad86bd1611447453c22c6fc"},
			{"Size", "4"},
		},
	},
	{
		name:  "curl get",
		hType: REQUEST,
		raw: "GET /test HTTP/1.1\r\n" +
			"User-Agent: curl/7.18.0 (i486-pc-linux-gnu) libcurl/7.18.0 OpenSSL/0.9.8g zlib/1.2.3.3 libidn/1.1\r\n" +
			"Host: 0.0.0.0=5000\r\n" +
			"Accept: */*\r\n" +
			"\r\n",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		method:                  GET,
		requestURL:              "/test",
		contentLength:           unused,
		headers: [][2]string{
			{"User-Agent", "curl/7.18.0 (i486-pc-linux-gnu) libcurl/7.18.0 OpenSSL/0.9.8g zlib/1.2.3.3 libidn/1.1"},
			{"Host", "0.0.0.0=5000"},
			{"Accept", "*/*"},
		},
	},
	{
		name:                    "firefox get",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET /favicon.ico HTTP/1.1\r\n" +
			"Host: 0.0.0.0=5000\r\n" +
			"User-Agent: Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9) Gecko/2008061015 Firefox/3.0\r\n" +
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n" +
			"Accept-Language: en-us,en;q=0.5\r\n" +
			"Accept-Encoding: gzip,deflate\r\n" +
			"Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7\r\n" +
			"Keep-Alive: 300\r\n" +
			"Connection: keep-alive\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "/favicon.ico",
		contentLength:        unused,
		headers: [][2]string{
			{"Host", "0.0.0.0=5000"},
			{"User-Agent", "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9) Gecko/2008061015 Firefox/3.0"},
			{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
			{"Accept-Language", "en-us,en;q=0.5"},
			{"Accept-Encoding", "gzip,deflate"},
			{"Accept-Charset", "ISO-8859-1,utf-8;q=0.7,*;q=0.7"},
			{"Keep-Alive", "300"},
			{"Connection", "keep-alive"},
		},
	},
	{
		name:                    "dumbluck",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET /dumbluck HTTP/1.1\r\n" +
			"aaaaaaaaaaaaa:++++++++++\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "/dumbluck",
		contentLength:        unused,
		headers: [][2]string{
			{"aaaaaaaaaaaaa", "++++++++++"},
		},
	},
	{
		name:                    "fragment in url",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET /forums/1/topics/2375?page=1#posts-17408 HTTP/1.1\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "/forums/1/topics/2375?page=1#posts-17408",
		contentLength:        unused,
	},
	{
		name:                    "get no headers no body",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET /get_no_headers_no_body/world HTTP/1.1\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "/get_no_headers_no_body/world",
		contentLength:        unused,
	},
	{
		name:                    "get one header no body",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET /get_one_header_no_body HTTP/1.1\r\n" +
			"Accept: */*\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "/get_one_header_no_body",
		contentLength:        unused,
		headers: [][2]string{
			{"Accept", "*/*"},
		},
	},
	{
		name:                    "get funky content length body hello",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET /get_funky_content_length_body_hello HTTP/1.0\r\n" +
			"conTENT-Length: 5\r\n" +
			"\r\n" +
			"HELLO",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            0,
		method:               GET,
		requestURL:           "/get_funky_content_length_body_hello",
		contentLength:        unused,
		headers: [][2]string{
			{"conTENT-Length", "5"},
		},
		body: "HELLO",
	},
	{
		name:                    "post identity body world",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "POST /post_identity_body_world?q=search#hey HTTP/1.1\r\n" +
			"Accept: */*\r\n" +
			"Content-Length: 5\r\n" +
			"\r\n" +
			"World",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               POST,
		requestURL:           "/post_identity_body_world?q=search#hey",
		contentLength:        unused,
		headers: [][2]string{
			{"Accept", "*/*"},
			{"Content-Length", "5"},
		},
		body: "World",
	},
	{
		name:                    "post - chunked body: all your base are belong to us",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "POST /post_chunked_all_your_base HTTP/1.1\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"\r\n" +
			"1e\r\nall your base are belong to us\r\n" +
			"0\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               POST,
		requestURL:           "/post_chunked_all_your_base",
		contentLength:        unused,
		headers: [][2]string{
			{"Transfer-Encoding", "chunked"},
		},
		body: "all your base are belong to us",
	},
	{
		name:                    "two chunks ; triple zero ending",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "POST /two_chunks_mult_zero_end HTTP/1.1\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"\r\n" +
			"5\r\nhello\r\n" +
			"6\r\n world\r\n" +
			"000\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               POST,
		requestURL:           "/two_chunks_mult_zero_end",
		contentLength:        unused,
		headers: [][2]string{
			{"Transfer-Encoding", "chunked"},
		},
		body: "hello world",
	},
	{
		name:                    "chunked with trailing headers. blech.",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "POST /chunked_w_trailing_headers HTTP/1.1\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"\r\n" +
			"5\r\nhello\r\n" +
			"6\r\n world\r\n" +
			"0\r\n" +
			"Vary: *\r\n" +
			"Content-Type: text/plain\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               POST,
		requestURL:           "/chunked_w_trailing_headers",
		contentLength:        unused,
		headers: [][2]string{
			{"Transfer-Encoding", "chunked"},
			{"Vary", "*"},
			{"Content-Type", "text/plain"},
		},
		body: "hello world",
	},
	{
		name:                    "with nonsense after the length",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "POST /chunked_w_nonsense_after_length HTTP/1.1\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"\r\n" +
			"5; ilovew3;whattheluck=aretheseparametersfor\r\nhello\r\n" +
			"6; blahblah; blah\r\n world\r\n" +
			"0\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               POST,
		requestURL:           "/chunked_w_nonsense_after_length",
		contentLength:        unused,
		headers: [][2]string{
			{"Transfer-Encoding", "chunked"},
		},
		body: "hello world",
	},
	{
		name:                    "with quotes",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw:                     "GET /with_\"stupid\"_quotes?foo=\"bar\" HTTP/1.1\r\n\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "/with_\"stupid\"_quotes?foo=\"bar\"",
		contentLength:        unused,
	},
	{
		name:                    "apachebench get",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET /test HTTP/1.0\r\n" +
			"Host: 0.0.0.0:5000\r\n" +
			"User-Agent: ApacheBench/2.3\r\n" +
			"Accept: */*\r\n\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            0,
		method:               GET,
		requestURL:           "/test",
		contentLength:        unused,
		headers: [][2]string{
			{"Host", "0.0.0.0:5000"},
			{"User-Agent", "ApacheBench/2.3"},
			{"Accept", "*/*"},
		},
	},
	{
		name:                    "query url with question mark",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw:                     "GET /test.cgi?foo=bar?baz HTTP/1.1\r\n\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "/test.cgi?foo=bar?baz",
		contentLength:        unused,
	},
	/* Some clients, especially after a POST in a keep-alive connection,
	 * will send an extra CRLF before the next request
	 */
	{
		name:                    "newline prefix get",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw:                     "\r\nGET /test HTTP/1.1\r\n\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "/test",
		contentLength:        unused,
	},
	{
		name:                    "upgrade request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET /demo HTTP/1.1\r\n" +
			"Host: example.com\r\n" +
			"Connection: Upgrade\r\n" +
			"Sec-WebSocket-Key2: 12998 5 Y3 1  .P00\r\n" +
			"Sec-WebSocket-Protocol: sample\r\n" +
			"Upgrade: WebSocket\r\n" +
			"Sec-WebSocket-Key1: 4 @1  46546xW%0l 1 5\r\n" +
			"Origin: http://example.com\r\n" +
			"\r\n" +
			"Hot diggity dogg",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		upgrade:              "Hot diggity dogg",
		method:               GET,
		requestURL:           "/demo",
		contentLength:        unused,
		headers: [][2]string{{"Host", "example.com"},
			{"Connection", "Upgrade"},
			{"Sec-WebSocket-Key2", "12998 5 Y3 1  .P00"},
			{"Sec-WebSocket-Protocol", "sample"},
			{"Upgrade", "WebSocket"},
			{"Sec-WebSocket-Key1", "4 @1  46546xW%0l 1 5"},
			{"Origin", "http://example.com"},
		},
	},
	{
		name:                    "connect request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "CONNECT 0-home0.netscape.com:443 HTTP/1.0\r\n" +
			"User-agent: Mozilla/1.1N\r\n" +
			"Proxy-authorization: basic aGVsbG86d29ybGQ=\r\n" +
			"\r\n" +
			"some data\r\n" +
			"and yet even more data",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            0,
		upgrade:              "some data\r\nand yet even more data",
		method:               CONNECT,
		requestURL:           "0-home0.netscape.com:443",
		contentLength:        unused,
		headers: [][2]string{
			{"User-agent", "Mozilla/1.1N"},
			{"Proxy-authorization", "basic aGVsbG86d29ybGQ="},
		},
	},
	{
		name:                    "report request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "REPORT /test HTTP/1.1\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               REPORT,
		requestURL:           "/test",
		contentLength:        unused,
	},
	/*
		{
			name:                    "request with no http version",
			hType:                   REQUEST,
			messageCompleteCbCalled: true,
			raw: "GET / \r\n" +
				"\r\n",

			shouldKeepAlive:      true,
			messageCompleteOnEOF: false,
			httpMajor:            0,
			httpMinor:            9,
			//method: HTTP_POST,
			requestURL:    "/",
			contentLength: unused,
		},
	*/
	{
		name:                    "m-search request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "M-SEARCH * HTTP/1.1\r\n" +
			"HOST: 239.255.255.250:1900\r\n" +
			"MAN: \"ssdp:discover\"\r\n" +
			"ST: \"ssdp:all\"\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               MSEARCH,
		requestURL:           "*",
		contentLength:        unused,
		headers: [][2]string{
			{"HOST", "239.255.255.250:1900"},
			{"MAN", "\"ssdp:discover\""},
			{"ST", "\"ssdp:all\""},
		},
	},
	/*
		{
			name:                    "line folding in header value",
			hType:                   REQUEST,
			messageCompleteCbCalled: true,
			raw: "GET / HTTP/1.1\r\n" +
				"Line1:   abc\r\n" +
				"\tdef\r\n" +
				" ghi\r\n" +
				"\t\tjkl\r\n" +
				"  mno \r\n" +
				"\t \tqrs\r\n" +
				"Line2: \t line2\t\r\n" +
				"Line3:\r\n" +
				" line3\r\n" +
				"Line4: \r\n" +
				" \r\n" +
				"Connection:\r\n" +
				" close\r\n" +
				"\r\n",

			shouldKeepAlive:      true,
			messageCompleteOnEOF: false,
			httpMajor:            1,
			httpMinor:            1,
			method:               GET,
			requestURL:           "/",
			contentLength:        unused,
			headers: [][2]string{
				{"Line1", "abc\tdef ghi\t\tjkl  mno \t \tqrs"},
				{"Line2", "line2\t"},
				{"Line3", "line3"},
				{"Line4", ""},
				{"Connection", "close"},
			},
		},
	*/
	{
		name:                    "host terminated by a query string",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET http://hypnotoad.org?hail=all HTTP/1.1\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "http://hypnotoad.org?hail=all",
		contentLength:        unused,
	},
	{
		name:                    "host:port terminated by a query string",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET http://hypnotoad.org:1234?hail=all HTTP/1.1\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "http://hypnotoad.org:1234?hail=all",
		contentLength:        unused,
	},
	{
		name:                    "host:port terminated by a space",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET http://hypnotoad.org:1234 HTTP/1.1\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "http://hypnotoad.org:1234",
		contentLength:        unused,
	},
	{
		name:                    "PATCH request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "PATCH /file.txt HTTP/1.1\r\n" +
			"Host: www.example.com\r\n" +
			"Content-Type: application/example\r\n" +
			"If-Match: \"e0023aa4e\"\r\n" +
			"Content-Length: 10\r\n" +
			"\r\n" +
			"cccccccccc",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		body:                 "cccccccccc",
		method:               PATCH,
		requestURL:           "/file.txt",
		contentLength:        unused,
		headers: [][2]string{
			{"Host", "www.example.com"},
			{"Content-Type", "application/example"},
			{"If-Match", "\"e0023aa4e\""},
			{"Content-Length", "10"},
		},
	},
	{
		name:                    "connect caps request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "CONNECT HOME0.NETSCAPE.COM:443 HTTP/1.0\r\n" +
			"User-agent: Mozilla/1.1N\r\n" +
			"Proxy-authorization: basic aGVsbG86d29ybGQ=\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            0,
		method:               CONNECT,
		requestURL:           "HOME0.NETSCAPE.COM:443",
		contentLength:        unused,
		headers: [][2]string{
			{"User-agent", "Mozilla/1.1N"},
			{"Proxy-authorization", "basic aGVsbG86d29ybGQ="},
		},
	},
	{
		name:                    "utf-8 path request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET /δ¶/δt/pope?q=1#narf HTTP/1.1\r\n" +
			"Host: github.com\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "/δ¶/δt/pope?q=1#narf",
		contentLength:        unused,
		headers: [][2]string{
			{"Host", "github.com"},
		},
	},
	{
		name:                    "hostname underscore",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "CONNECT home_0.netscape.com:443 HTTP/1.0\r\n" +
			"User-agent: Mozilla/1.1N\r\n" +
			"Proxy-authorization: basic aGVsbG86d29ybGQ=\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            0,
		method:               CONNECT,
		requestURL:           "home_0.netscape.com:443",
		contentLength:        unused,
		headers: [][2]string{
			{"User-agent", "Mozilla/1.1N"},
			{"Proxy-authorization", "basic aGVsbG86d29ybGQ="},
		},
	},
	{
		name:                    "eat CRLF between requests, no \"Connection: close\" header",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "POST / HTTP/1.1\r\n" +
			"Host: www.example.com\r\n" +
			"Content-Type: application/x-www-form-urlencoded\r\n" +
			"Content-Length: 4\r\n" +
			"\r\n" +
			"q=42\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               POST,
		requestURL:           "/",
		body:                 "q=42",
		contentLength:        unused,
		headers: [][2]string{
			{"Host", "www.example.com"},
			{"Content-Type", "application/x-www-form-urlencoded"},
			{"Content-Length", "4"},
		},
	},
	{
		name:                    "PURGE request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "PURGE /file.txt HTTP/1.1\r\n" +
			"Host: www.example.com\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               PURGE,
		requestURL:           "/file.txt",
		contentLength:        unused,
		headers: [][2]string{
			{"Host", "www.example.com"},
		},
	},
	{
		name:                    "SEARCH request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "SEARCH / HTTP/1.1\r\n" +
			"Host: www.example.com\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               SEARCH,
		requestURL:           "/",
		contentLength:        unused,
		headers: [][2]string{
			{"Host", "www.example.com"},
		},
	},
	{
		name:                    "host:port and basic_auth",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET http://a%12:b!&*$@hypnotoad.org:1234/toto HTTP/1.1\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "http://a%12:b!&*$@hypnotoad.org:1234/toto",
		contentLength:        unused,
	},
	//35
	/*
		{
			name:                    "multiple connection header values with folding",
			hType:                   REQUEST,
			messageCompleteCbCalled: true,
			raw: "GET /demo HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"Connection: Something,\r\n" +
				" Upgrade, ,Keep-Alive\r\n" +
				"Sec-WebSocket-Key2: 12998 5 Y3 1  .P00\r\n" +
				"Sec-WebSocket-Protocol: sample\r\n" +
				"Upgrade: WebSocket\r\n" +
				"Sec-WebSocket-Key1: 4 @1  46546xW%0l 1 5\r\n" +
				"Origin: http://example.com\r\n" +
				"\r\n" +
				"Hot diggity dogg",

			shouldKeepAlive:      true,
			messageCompleteOnEOF: false,
			httpMajor:            1,
			httpMinor:            1,
			method:               PATCH,
			requestURL:           "/demo",
			contentLength:        unused,
			upgrade:              "Hot diggity dogg",
			headers: [][2]string{
				{"Host", "example.com"},
				{"Connection", "Something, Upgrade, ,Keep-Alive"},
				{"Sec-WebSocket-Key2", "12998 5 Y3 1  .P00"},
				{"Sec-WebSocket-Protocol", "sample"},
				{"Upgrade", "WebSocket"},
				{"Sec-WebSocket-Key1", "4 @1  46546xW%0l 1 5"},
				{"Origin", "http://example.com"},
			},
		},
	*/
	// 36
	{
		name:                    "multiple connection header values with folding and lws",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET /demo HTTP/1.1\r\n" +
			"Connection: keep-alive, upgrade\r\n" +
			"Upgrade: WebSocket\r\n" +
			"\r\n" +
			"Hot diggity dogg",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               GET,
		requestURL:           "/demo",
		contentLength:        unused,
		upgrade:              "Hot diggity dogg",
		headers: [][2]string{
			{"Connection", "keep-alive, upgrade"},
			{"Upgrade", "WebSocket"},
		},
	},
	// 37
	/*
		{
			name:                    "multiple connection header values with folding and lws",
			hType:                   REQUEST,
			messageCompleteCbCalled: true,
			raw: "GET /demo HTTP/1.1\r\n" +
				"Connection: keep-alive, \r\n upgrade\r\n" +
				"Upgrade: WebSocket\r\n" +
				"\r\n" +
				"Hot diggity dogg",

			shouldKeepAlive:      true,
			messageCompleteOnEOF: false,
			httpMajor:            1,
			httpMinor:            1,
			method:               PATCH,
			requestURL:           "/demo",
			upgrade:              "Hot diggity dogg",
			headers: [][2]string{
				{"Connection", "keep-alive,  upgrade"},
				{"Upgrade", "WebSocket"},
			},
			contentLength: unused,
		},
	*/
	// 38
	{
		name:                    "upgrade post request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "POST /demo HTTP/1.1\r\n" +
			"Host: example.com\r\n" +
			"Connection: Upgrade\r\n" +
			"Upgrade: HTTP/2.0\r\n" +
			"Content-Length: 15\r\n" +
			"\r\n" +
			"sweet post body" +
			"Hot diggity dogg",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               POST,
		requestURL:           "/demo",
		upgrade:              "Hot diggity dogg",
		body:                 "sweet post body",
		headers: [][2]string{
			{"Host", "example.com"},
			{"Connection", "Upgrade"},
			{"Upgrade", "HTTP/2.0"},
			{"Content-Length", "15"},
		},
		contentLength: unused,
	},
	// 39
	{
		name:                    "connect with body request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "CONNECT foo.bar.com:443 HTTP/1.0\r\n" +
			"User-agent: Mozilla/1.1N\r\n" +
			"Proxy-authorization: basic aGVsbG86d29ybGQ=\r\n" +
			"Content-Length: 10\r\n" +
			"\r\n" +
			"blarfcicle",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            0,
		method:               CONNECT,
		requestURL:           "foo.bar.com:443",
		contentLength:        unused,
		upgrade:              "blarfcicle",
		headers: [][2]string{
			{"User-agent", "Mozilla/1.1N"},
			{"Proxy-authorization", "basic aGVsbG86d29ybGQ="},
			{"Content-Length", "10"},
		},
	},
	// 40
	{
		name:                    "link request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "LINK /images/my_dog.jpg HTTP/1.1\r\n" +
			"Host: example.com\r\n" +
			"Link: <http://example.com/profiles/joe>; rel=\"tag\"\r\n" +
			"Link: <http://example.com/profiles/sally>; rel=\"tag\"\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               LINK,
		requestURL:           "/images/my_dog.jpg",
		contentLength:        unused,
		headers: [][2]string{
			{"Host", "example.com"},
			{"Link", "<http://example.com/profiles/joe>; rel=\"tag\""},
			{"Link", "<http://example.com/profiles/sally>; rel=\"tag\""},
		},
	},
	// 41
	{
		name:                    "unlink request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "UNLINK /images/my_dog.jpg HTTP/1.1\r\n" +
			"Host: example.com\r\n" +
			"Link: <http://example.com/profiles/sally>; rel=\"tag\"\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               UNLINK,
		requestURL:           "/images/my_dog.jpg",
		contentLength:        unused,
		headers: [][2]string{
			{"Host", "example.com"},
			{"Link", "<http://example.com/profiles/sally>; rel=\"tag\""},
		},
	},
	// 42
	{
		name:                    "source request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "SOURCE /music/sweet/music HTTP/1.1\r\n" +
			"Host: example.com\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               SOURCE,
		requestURL:           "/music/sweet/music",
		contentLength:        unused,
		headers: [][2]string{
			{"Host", "example.com"},
		},
	},
	// 43
	{
		name:                    "source request",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "SOURCE /music/sweet/music ICE/1.0\r\n" +
			"Host: example.com\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            0,
		method:               SOURCE,
		requestURL:           "/music/sweet/music",
		contentLength:        unused,
		headers: [][2]string{
			{"Host", "example.com"},
		},
	},
	// 44
	{
		name:                    "post - multi coding transfer-encoding chunked body",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "POST / HTTP/1.1\r\n" +
			"Transfer-Encoding: deflate, chunked\r\n" +
			"\r\n" +
			"1e\r\nall your base are belong to us\r\n" +
			"0\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEOF: false,
		httpMajor:            1,
		httpMinor:            1,
		method:               POST,
		requestURL:           "/",
		contentLength:        unused,
		body:                 "all your base are belong to us",
		headers: [][2]string{
			{"Transfer-Encoding", "deflate, chunked"},
		},
	},
	/*
		{
			name:                    "post - multi line coding transfer-encoding chunked body",
			hType:                   REQUEST,
			messageCompleteCbCalled: true,
			raw: "POST / HTTP/1.1\r\n" +
				"Transfer-Encoding: deflate,\r\n" +
				" chunked\r\n" +
				"\r\n" +
				"1e\r\nall your base are belong to us\r\n" +
				"0\r\n" +
				"\r\n",

			shouldKeepAlive:      true,
			messageCompleteOnEOF: false,
			httpMajor:            1,
			httpMinor:            1,
			method:               POST,
			requestURL:           "/",
			contentLength:        unused,
			headers: [][2]string{
				{"Transfer-Encoding", "deflate, chunked"},
			},
			body: "all your base are belong to us",
		},
	*/
	/*
		{
			name:                    "chunked with content-length set, allow_chunked_length flag is set",
			hType:                   REQUEST,
			messageCompleteCbCalled: true,
			raw: "POST /chunked_w_content_length HTTP/1.1\r\n" +
				"Content-Length: 10\r\n" +
				"Transfer-Encoding: chunked\r\n" +
				"\r\n" +
				"5; ilovew3;whattheluck=aretheseparametersfor\r\nhello\r\n" +
				"6; blahblah; blah\r\n world\r\n" +
				"0\r\n" +
				"\r\n",

			shouldKeepAlive:      true,
			messageCompleteOnEOF: false,
			httpMajor:            1,
			httpMinor:            1,
			method:               PATCH,
			requestURL:           "/chunked_w_content_length",
			contentLength:        unused,
			headers: [][2]string{
				{"Content-Length", "10"},
				{"Transfer-Encoding", "chunked"},
			},
			body: "hello world",
		},
	*/
}

//var requestsDebug = []message{}

var responses = []message{
	{
		name:  "google 301",
		hType: RESPONSE,
		raw: "HTTP/1.1 301 Moved Permanently\r\n" +
			"Location: http://www.google.com/\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"Date: Sun, 26 Apr 2009 11:11:49 GMT\r\n" +
			"Expires: Tue, 26 May 2009 11:11:49 GMT\r\n" +
			"X-$PrototypeBI-Version: 1.6.0.3\r\n" + /* $ char in header field */
			"Cache-Control: public, max-age=2592000\r\n" +
			"Server: gws\r\n" +
			"Content-Length:  219  \r\n" +
			"\r\n" +
			"<HTML><HEAD><meta http-equiv=\"content-type\" content=\"text/html;charset=utf-8\">\n" +
			"<TITLE>301 Moved</TITLE></HEAD><BODY>\n" +
			"<H1>301 Moved</H1>\n" +
			"The document has moved\n" +
			"<A HREF=\"http://www.google.com/\">here</A>.\r\n" +
			"</BODY></HTML>\r\n",

		statusCode:              301,
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		responseStatus:          "Moved Permanently",
		//method: HTTP_GET,
		headers: [][2]string{
			{"Location", "http://www.google.com/"},
			{"Content-Type", "text/html; charset=UTF-8"},
			{"Date", "Sun, 26 Apr 2009 11:11:49 GMT"},
			{"Expires", "Tue, 26 May 2009 11:11:49 GMT"},
			{"X-$PrototypeBI-Version", "1.6.0.3"},
			{"Cache-Control", "public, max-age=2592000"},
			{"Server", "gws"},
			{"Content-Length", " 219  "},
		},
		contentLength: unused,
		body: "<HTML><HEAD><meta http-equiv=\"content-type\" content=\"text/html;charset=utf-8\">\n" +
			"<TITLE>301 Moved</TITLE></HEAD><BODY>\n" +
			"<H1>301 Moved</H1>\n" +
			"The document has moved\n" +
			"<A HREF=\"http://www.google.com/\">here</A>.\r\n" +
			"</BODY></HTML>\r\n",
	},
	{
		name:  "no content-length response",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\r\n" +
			"Date: Tue, 04 Aug 2009 07:59:32 GMT\r\n" +
			"Server: Apache\r\n" +
			"X-Powered-By: Servlet/2.5 JSP/2.1\r\n" +
			"Content-Type: text/xml; charset=utf-8\r\n" +
			"Connection: close\r\n" +
			"\r\n" +
			"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
			"<SOAP-ENV:Envelope xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\">\n" +
			"  <SOAP-ENV:Body>\n" +
			"    <SOAP-ENV:Fault>\n" +
			"       <faultcode>SOAP-ENV:Client</faultcode>\n" +
			"       <faultstring>Client Error</faultstring>\n" +
			"    </SOAP-ENV:Fault>\n" +
			"  </SOAP-ENV:Body>\n" +
			"</SOAP-ENV:Envelope>",
		statusCode:              200,
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		responseStatus:          "OK",
		//method: HTTP_GET,
		headers: [][2]string{
			{"Date", "Tue, 04 Aug 2009 07:59:32 GMT"},
			{"Server", "Apache"},
			{"X-Powered-By", "Servlet/2.5 JSP/2.1"},
			{"Content-Type", "text/xml; charset=utf-8"},
			{"Connection", "close"},
		},
		contentLength: unused,
		body: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
			"<SOAP-ENV:Envelope xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\">\n" +
			"  <SOAP-ENV:Body>\n" +
			"    <SOAP-ENV:Fault>\n" +
			"       <faultcode>SOAP-ENV:Client</faultcode>\n" +
			"       <faultstring>Client Error</faultstring>\n" +
			"    </SOAP-ENV:Fault>\n" +
			"  </SOAP-ENV:Body>\n" +
			"</SOAP-ENV:Envelope>",
	},
	{
		name:                    "404 no headers no body",
		hType:                   RESPONSE,
		raw:                     "HTTP/1.1 404 Not Found\r\n\r\n",
		statusCode:              404,
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		responseStatus:          "Not Found",
		//method: HTTP_GET,
		contentLength: unused,
	},
	{
		name:                    "301 no response phrase",
		hType:                   RESPONSE,
		raw:                     "HTTP/1.1 301\r\n\r\n",
		statusCode:              301,
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		//method: HTTP_GET,
		contentLength: unused,
	},
	{
		name:  "200 trailing space on chunked body",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/plain\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"\r\n" +
			"25  \r\n" +
			"This is the data in the first chunk\r\n" +
			"\r\n" +
			"1C\r\n" +
			"and this is the second one\r\n" +
			"\r\n" +
			"0  \r\n" +
			"\r\n",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		//method: HTTP_GET,
		contentLength: unused,
		body: "This is the data in the first chunk\r\n" +
			"and this is the second one\r\n",
		headers: [][2]string{
			{"Content-Type", "text/plain"},
			{"Transfer-Encoding", "chunked"},
		},
	},
	{
		name:  "no carriage ret",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\n" +
			"Content-Type: text/html; charset=utf-8\n" +
			"Connection: close\n" +
			"\n" +
			"these headers are from http://news.ycombinator.com/",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		//method: HTTP_GET,
		contentLength: unused,
		body:          "these headers are from http://news.ycombinator.com/",
		headers: [][2]string{
			{"Content-Type", "text/html; charset=utf-8"},
			{"Connection", "close"},
		},
	},
	{
		name:  "proxy connection",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"Content-Length: 11\r\n" +
			"Proxy-Connection: close\r\n" +
			"Date: Thu, 31 Dec 2009 20:55:48 +0000\r\n" +
			"\r\n" +
			"hello world",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		//method: HTTP_GET,
		contentLength: unused,
		body:          "hello world",
		headers: [][2]string{
			{"Content-Type", "text/html; charset=UTF-8"},
			{"Content-Length", "11"},
			{"Proxy-Connection", "close"},
			{"Date", "Thu, 31 Dec 2009 20:55:48 +0000"},
		},
	},
	{
		name:  "underscore header key",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\r\n" +
			"Server: DCLK-AdSvr\r\n" +
			"Content-Type: text/xml\r\n" +
			"Content-Length: 0\r\n" +
			"DCLK_imp: v7;x;114750856;0-0;0;17820020;0/0;21603567/21621457/1;;~okv=;dcmt=text/xml;;~cs=o\r\n\r\n",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Server", "DCLK-AdSvr"},
			{"Content-Type", "text/xml"},
			{"Content-Length", "0"},
			{"DCLK_imp", "v7;x;114750856;0-0;0;17820020;0/0;21603567/21621457/1;;~okv=;dcmt=text/xml;;~cs=o"},
		},
	},
	{
		name:  "bonjourmadame.fr",
		hType: RESPONSE,
		raw: "HTTP/1.0 301 Moved Permanently\r\n" +
			"Date: Thu, 03 Jun 2010 09:56:32 GMT\r\n" +
			"Server: Apache/2.2.3 (Red Hat)\r\n" +
			"Cache-Control: public\r\n" +
			"Pragma: \r\n" +
			"Location: http://www.bonjourmadame.fr/\r\n" +
			"Vary: Accept-Encoding\r\n" +
			"Content-Length: 0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"Connection: keep-alive\r\n" +
			"\r\n",
		statusCode:              301,
		responseStatus:          "Moved Permanently",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               0,
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Date", "Thu, 03 Jun 2010 09:56:32 GMT"},
			{"Server", "Apache/2.2.3 (Red Hat)"},
			{"Cache-Control", "public"},
			{"Pragma", ""},
			{"Location", "http://www.bonjourmadame.fr/"},
			{"Vary", "Accept-Encoding"},
			{"Content-Length", "0"},
			{"Content-Type", "text/html; charset=UTF-8"},
			{"Connection", "keep-alive"},
		},
	},
	{
		name:  "field underscore",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\r\n" +
			"Date: Tue, 28 Sep 2010 01:14:13 GMT\r\n" +
			"Server: Apache\r\n" +
			"Cache-Control: no-cache, must-revalidate\r\n" +
			"Expires: Mon, 26 Jul 1997 05:00:00 GMT\r\n" +
			".et-Cookie: PlaxoCS=1274804622353690521; path=/; domain=.plaxo.com\r\n" +
			"Vary: Accept-Encoding\r\n" +
			"_eep-Alive: timeout=45\r\n" + /* semantic value ignored */
			"_onnection: Keep-Alive\r\n" + /* semantic value ignored */
			"Transfer-Encoding: chunked\r\n" +
			"Content-Type: text/html\r\n" +
			"Connection: close\r\n" +
			"\r\n" +
			"0\r\n\r\n",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Date", "Tue, 28 Sep 2010 01:14:13 GMT"},
			{"Server", "Apache"},
			{"Cache-Control", "no-cache, must-revalidate"},
			{"Expires", "Mon, 26 Jul 1997 05:00:00 GMT"},
			{".et-Cookie", "PlaxoCS=1274804622353690521; path=/; domain=.plaxo.com"},
			{"Vary", "Accept-Encoding"},
			{"_eep-Alive", "timeout=45"},
			{"_onnection", "Keep-Alive"},
			{"Transfer-Encoding", "chunked"},
			{"Content-Type", "text/html"},
			{"Connection", "close"},
		},
	},
	{
		name:  "non-ASCII in status line",
		hType: RESPONSE,
		raw: "HTTP/1.1 500 Oriëntatieprobleem\r\n" +
			"Date: Fri, 5 Nov 2010 23:07:12 GMT+2\r\n" +
			"Content-Length: 0\r\n" +
			"Connection: close\r\n" +
			"\r\n",
		statusCode:              500,
		responseStatus:          "Oriëntatieprobleem",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Date", "Fri, 5 Nov 2010 23:07:12 GMT+2"},
			{"Content-Length", "0"},
			{"Connection", "close"},
		},
	},
	{
		name:  "http version 0.9",
		hType: RESPONSE,
		raw: "HTTP/0.9 200 OK\r\n" +
			"\r\n",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               0,
		httpMinor:               9,
		//method: HTTP_GET,
		contentLength: unused,
	},
	{
		name:  "neither content-length nor transfer-encoding response",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/plain\r\n" +
			"\r\n" +
			"hello world",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		body:                    "hello world",
		headers: [][2]string{
			{"Content-Type", "text/plain"},
		},
		//method: HTTP_GET,
		contentLength: unused,
	},
	{
		name:  "HTTP/1.0 with keep-alive and EOF-terminated 200 status",
		hType: RESPONSE,
		raw: "HTTP/1.0 200 OK\r\n" +
			"Connection: keep-alive\r\n" +
			"\r\n",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               0,
		headers: [][2]string{
			{"Connection", "keep-alive"},
		},
		//method: HTTP_GET,
		contentLength: unused,
	},
	{
		name:  "HTTP/1.0 with keep-alive and a 204 status",
		hType: RESPONSE,
		raw: "HTTP/1.0 204 No content\r\n" +
			"Connection: keep-alive\r\n" +
			"\r\n",
		statusCode:              204,
		responseStatus:          "No content",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               0,
		headers: [][2]string{
			{"Connection", "keep-alive"},
		},
		//method: HTTP_GET,
		contentLength: unused,
	},
	{
		name:  "HTTP/1.1 with an EOF-terminated 200 status",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\r\n" +
			"\r\n",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		//method: HTTP_GET,
		contentLength: unused,
	},
	{
		name:  "HTTP/1.1 with a 204 status",
		hType: RESPONSE,
		raw: "HTTP/1.1 204 No content\r\n" +
			"\r\n",
		statusCode:              204,
		responseStatus:          "No content",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		//method: HTTP_GET,
		contentLength: unused,
	},
	{
		name:  "HTTP/1.1 with a 204 status and keep-alive disabled",
		hType: RESPONSE,
		raw: "HTTP/1.1 204 No content\r\n" +
			"Connection: close\r\n" +
			"\r\n",
		statusCode:              204,
		responseStatus:          "No content",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Connection", "close"},
		},
	},
	{
		name:  "field space",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\r\n" +
			"Server: Microsoft-IIS/6.0\r\n" +
			"X-Powered-By: ASP.NET\r\n" +
			"en-US Content-Type: text/xml\r\n" + /* this is the problem */
			"Content-Type: text/xml\r\n" +
			"Content-Length: 16\r\n" +
			"Date: Fri, 23 Jul 2010 18:45:38 GMT\r\n" +
			"Connection: keep-alive\r\n" +
			"\r\n" +
			"<xml>hello</xml>",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		body:                    "<xml>hello</xml>",
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Server", "Microsoft-IIS/6.0"},
			{"X-Powered-By", "ASP.NET"},
			{"en-US Content-Type", "text/xml"},
			{"Content-Type", "text/xml"},
			{"Content-Length", "16"},
			{"Date", "Fri, 23 Jul 2010 18:45:38 GMT"},
			{"Connection", "keep-alive"},
		},
	},
	{
		name:  "amazon.com",
		hType: RESPONSE,
		raw: "HTTP/1.1 301 MovedPermanently\r\n" +
			"Date: Wed, 15 May 2013 17:06:33 GMT\r\n" +
			"Server: Server\r\n" +
			"x-amz-id-1: 0GPHKXSJQ826RK7GZEB2\r\n" +
			"p3p: policyref=\"http://www.amazon.com/w3c/p3p.xml\",CP=\"CAO DSP LAW CUR ADM IVAo IVDo CONo OTPo OUR DELi PUBi OTRi BUS PHY ONL UNI PUR FIN COM NAV INT DEM CNT STA HEA PRE LOC GOV OTC \"\r\n" +
			"x-amz-id-2: STN69VZxIFSz9YJLbz1GDbxpbjG6Qjmmq5E3DxRhOUw+Et0p4hr7c/Q8qNcx4oAD\r\n" +
			"Location: http://www.amazon.com/Dan-Brown/e/B000AP9DSU/ref=s9_pop_gw_al1?_encoding=UTF8&refinementId=618073011&pf_rd_m=ATVPDKIKX0DER&pf_rd_s=center-2&pf_rd_r=0SHYY5BZXN3KR20BNFAY&pf_rd_t=101&pf_rd_p=1263340922&pf_rd_i=507846\r\n" +
			"Vary: Accept-Encoding,User-Agent\r\n" +
			"Content-Type: text/html; charset=ISO-8859-1\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"\r\n" +
			"1\r\n" +
			"\n\r\n" +
			"0\r\n" +
			"\r\n",
		statusCode:              301,
		responseStatus:          "MovedPermanently",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		body:                    "\n",
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Date", "Wed, 15 May 2013 17:06:33 GMT"},
			{"Server", "Server"},
			{"x-amz-id-1", "0GPHKXSJQ826RK7GZEB2"},
			{"p3p", "policyref=\"http://www.amazon.com/w3c/p3p.xml\",CP=\"CAO DSP LAW CUR ADM IVAo IVDo CONo OTPo OUR DELi PUBi OTRi BUS PHY ONL UNI PUR FIN COM NAV INT DEM CNT STA HEA PRE LOC GOV OTC \""},
			{"x-amz-id-2", "STN69VZxIFSz9YJLbz1GDbxpbjG6Qjmmq5E3DxRhOUw+Et0p4hr7c/Q8qNcx4oAD"},
			{"Location", "http://www.amazon.com/Dan-Brown/e/B000AP9DSU/ref=s9_pop_gw_al1?_encoding=UTF8&refinementId=618073011&pf_rd_m=ATVPDKIKX0DER&pf_rd_s=center-2&pf_rd_r=0SHYY5BZXN3KR20BNFAY&pf_rd_t=101&pf_rd_p=1263340922&pf_rd_i=507846"},
			{"Vary", "Accept-Encoding,User-Agent"},
			{"Content-Type", "text/html; charset=ISO-8859-1"},
			{"Transfer-Encoding", "chunked"},
		},
	},
	{
		name:  "empty reason phrase after space",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 \r\n" +
			"\r\n",
		statusCode:              200,
		responseStatus:          "",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		//method: HTTP_GET,
		contentLength: unused,
	},
	{
		name:  "Content-Length-X",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\r\n" +
			"Content-Length-X: 0\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"\r\n" +
			"2\r\n" +
			"OK\r\n" +
			"0\r\n" +
			"\r\n",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		body:                    "OK",
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Content-Length-X", "0"},
			{"Transfer-Encoding", "chunked"},
		},
	},
	{
		name:  "HTTP 101 response with Upgrade header",
		hType: RESPONSE,
		raw: "HTTP/1.1 101 Switching Protocols\r\n" +
			"Connection: upgrade\r\n" +
			"Upgrade: h2c\r\n" +
			"\r\n" +
			"proto",
		statusCode:              101,
		responseStatus:          "Switching Protocols",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		upgrade:                 "proto",
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Connection", "upgrade"},
			{"Upgrade", "h2c"},
		},
	},
	{
		name:  "HTTP 101 response with Upgrade and Content-Length header",
		hType: RESPONSE,
		raw: "HTTP/1.1 101 Switching Protocols\r\n" +
			"Connection: upgrade\r\n" +
			"Upgrade: h2c\r\n" +
			"Content-Length: 4\r\n" +
			"\r\n" +
			"body" +
			"proto",
		statusCode:              101,
		responseStatus:          "Switching Protocols",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		upgrade:                 "proto",
		body:                    "body",
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Connection", "upgrade"},
			{"Upgrade", "h2c"},
			{"Content-Length", "4"},
		},
	},
	{
		name:  "HTTP 101 response with Upgrade and Transfer-Encoding header",
		hType: RESPONSE,
		raw: "HTTP/1.1 101 Switching Protocols\r\n" +
			"Connection: upgrade\r\n" +
			"Upgrade: h2c\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"\r\n" +
			"2\r\n" +
			"bo\r\n" +
			"2\r\n" +
			"dy\r\n" +
			"0\r\n" +
			"\r\n" +
			"proto",
		statusCode:              101,
		responseStatus:          "Switching Protocols",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		upgrade:                 "proto",
		body:                    "body",
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Connection", "upgrade"},
			{"Upgrade", "h2c"},
			{"Transfer-Encoding", "chunked"},
		},
	},
	{
		name:  "HTTP 200 response with Upgrade header",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\r\n" +
			"Connection: upgrade\r\n" +
			"Upgrade: h2c\r\n" +
			"\r\n" +
			"body",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		body:                    "body",
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Connection", "upgrade"},
			{"Upgrade", "h2c"},
		},
	},
	{
		name:  "HTTP 200 response with Upgrade and Content-Length header",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\r\n" +
			"Connection: upgrade\r\n" +
			"Upgrade: h2c\r\n" +
			"Content-Length: 4\r\n" +
			"\r\n" +
			"body",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		body:                    "body",
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Connection", "upgrade"},
			{"Upgrade", "h2c"},
			{"Content-Length", "4"},
		},
	},
	{
		name:  "HTTP 200 response with Upgrade and Transfer-Encoding header",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\r\n" +
			"Connection: upgrade\r\n" +
			"Upgrade: h2c\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"\r\n" +
			"2\r\n" +
			"bo\r\n" +
			"2\r\n" +
			"dy\r\n" +
			"0\r\n" +
			"\r\n",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEOF:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		body:                    "body",
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Connection", "upgrade"},
			{"Upgrade", "h2c"},
			{"Transfer-Encoding", "chunked"},
		},
	},
}

var settingTest Setting = Setting{
	MessageBegin: func(p *Parser, _ int) {
		m := p.GetUserData().(*message)
		m.messageBeginCbCalled = true
		m.hType = p.hType
	},
	URL: func(p *Parser, url []byte, _ int) {
		m := p.GetUserData().(*message)
		m.requestURL += string(url)
	},
	Status: func(p *Parser, status []byte, _ int) {
		m := p.GetUserData().(*message)
		m.responseStatus += string(status)
	},
	HeaderField: func(p *Parser, headerField []byte, _ int) {
		m := p.GetUserData().(*message)
		m.headers = append(m.headers, [2]string{string(headerField), ""})
	},
	HeaderValue: func(p *Parser, headerValue []byte, _ int) {
		m := p.GetUserData().(*message)
		m.headers[len(m.headers)-1][1] = string(headerValue)
	},
	HeadersComplete: func(p *Parser, _ int) {
		m := p.GetUserData().(*message)
		m.headersCompleteCbCalled = true
	},
	Body: func(p *Parser, body []byte, _ int) {
		m := p.GetUserData().(*message)
		m.body += string(body)
	},
	MessageComplete: func(p *Parser, _ int) {
		m := p.GetUserData().(*message)
		m.method = p.Method
		m.messageCompleteCbCalled = true
		m.statusCode = int(p.StatusCode)
		m.httpMajor = uint16(p.Major)
		m.httpMinor = uint16(p.Minor)

	},
}

func parse(p *Parser, data string) (int, error) {
	return p.Execute(&settingTest, []byte(data))
}

func testMessage(t *testing.T, m *message) {
	for msg1len := 0; msg1len < len(m.raw); msg1len++ {
		p := New(m.hType)
		got := &message{}
		p.SetUserData(got)

		msg1Message := m.raw[:msg1len]
		msg2Message := m.raw[msg1len:]

		var (
			n1   int
			err1 error
			err  error
			n2   int
			data string
		)

		// 模拟第1次读包
		if msg1len > 0 {
			n1, err1 = parse(p, msg1Message)
			if err1 != nil {
				t.Errorf("msg1len:%d, msg1(%s)", msg1len, msg1Message)
				return
			}
			// 如果有upgrade状态, 就不需要再重复送往数据
			if p.ReadyUpgradeData() {
				//if p.callMessageComplete && p.Upgrade {
				got.upgrade = msg1Message[n1:] + msg2Message
				msg1Message = ""
				goto test
			}

			msg1Message = msg1Message[n1:]
		}

		// 模拟第2次读包
		data = msg1Message + msg2Message
		n2, err = parse(p, data)
		if got.messageCompleteCbCalled && p.Upgrade {
			got.upgrade += data[n2:]
			goto test
		}
		if err != nil {
			t.Error(err)
			return
		}

	test:

		// flush 解析器
		_, err = parse(p, "")
		if err != nil {
			t.Error(err)
			return
		}
		if b := m.eq(t, got); !b {
			t.Logf("msg1.len:%d, msg2.len:%d, test case name:%s\n", len(msg1Message), len(msg2Message), m.name)
			t.Logf("msg1len:%d,  msg1(%s)", msg1len, msg1Message)
			t.Logf("          ,  msg2(%s)", msg2Message)
			t.Logf("upgrade:%t, got.messageCompleteCbCalled:%t, data:(%s)", p.Upgrade, got.messageCompleteCbCalled, data[n2:])
			break
		}

	}
}

func Test_Message(t *testing.T) {
	/*
		for _, req := range requestsDebug {
			//for _, req := range requests[len(requests)-1:] {
			testMessage(t, &req)
			_ = req
		}
	*/
	for _, req := range requests {
		testMessage(t, &req)
		_ = req
	}

	for _, rsp := range responses {
		testMessage(t, &rsp)
		_ = rsp
	}
}

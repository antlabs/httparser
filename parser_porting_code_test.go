// Copyright 2020 guonaihong. All rights reserved.
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
	"testing"

	"github.com/stretchr/testify/assert"
)

type message struct {
	name  string
	hType ReqOrRsp
	raw   string
	//httpMethod      method
	statusCode     int
	responseStatus string
	requestPath    string
	requestUrl     string
	body           string
	bodySize       string
	host           string
	userinfo       string
	port           uint16
	//enum { NONE=0, FIELD, VALUE } last_header_element; //TODO最近移值
	headers         [][2]string
	shouldKeepAlive bool

	numChunks         int
	numChunksComplete int
	chunkLengths      []int

	upgrade string // upgraded body

	httpMajor     uint16
	httpMinor     uint16
	contentLength int32

	messageBeginCbCalled    bool
	headersCompleteCbCalled bool
	messageCompleteCbCalled bool
	statusCbCalled          bool
	messageCompleteOnEof    bool
	bodyIsFinal             int
	allowChunkedLength      bool
}

func (m *message) eq(t *testing.T, m2 *message) bool {
	b := assert.Equal(t, m.messageCompleteCbCalled, m2.messageCompleteCbCalled, "messageCompleteCbCalled")
	if !b {
		return false
	}

	b = assert.Equal(t, m.headers, m2.headers)
	if !b {
		return false
	}

	b = assert.Equal(t, m.httpMajor, m2.httpMajor, "major")
	if !b {
		return false
	}
	b = assert.Equal(t, m.httpMinor, m2.httpMinor, "minor")
	if !b {
		return false
	}
	b = assert.Equal(t, m.hType, m2.hType, "htype")
	if !b {
		return false
	}

	b = assert.Equal(t, m.requestUrl, m2.requestUrl, "request url")
	if !b {
		return false
	}

	b = assert.Equal(t, m.body, m2.body, "body")
	if !b {
		return false
	}

	b = assert.Equal(t, m.responseStatus, m2.responseStatus, "responseStatus")
	if !b {
		return false
	}

	b = assert.Equal(t, m.upgrade, m2.upgrade, "upgrade")
	if !b {
		return false
	}
	return true
}

var requests = []message{
	{
		name:  "curl get",
		hType: REQUEST,
		raw: "GET /test HTTP/1.1\r\n" +
			"User-Agent: curl/7.18.0 (i486-pc-linux-gnu) libcurl/7.18.0 OpenSSL/0.9.8g zlib/1.2.3.3 libidn/1.1\r\n" +
			"Host: 0.0.0.0=5000\r\n" +
			"Accept: */*\r\n" +
			"\r\n",
		shouldKeepAlive:         true,
		messageCompleteOnEof:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		//method: HTTP_GET,
		requestUrl:    "/test",
		contentLength: unused,
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
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_GET,
		requestUrl:    "/favicon.ico",
		contentLength: unused,
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
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_GET,
		requestUrl:    "/dumbluck",
		contentLength: unused,
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
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_GET,
		requestUrl:    "/forums/1/topics/2375?page=1#posts-17408",
		contentLength: unused,
	},
	{
		name:                    "get no headers no body",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET /get_no_headers_no_body/world HTTP/1.1\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_GET,
		requestUrl:    "/get_no_headers_no_body/world",
		contentLength: unused,
	},
	{
		name:                    "get one header no body",
		hType:                   REQUEST,
		messageCompleteCbCalled: true,
		raw: "GET /get_one_header_no_body HTTP/1.1\r\n" +
			"Accept: */*\r\n" +
			"\r\n",

		shouldKeepAlive:      true,
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_GET,
		requestUrl:    "/get_one_header_no_body",
		contentLength: unused,
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
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            0,
		//method: HTTP_GET,
		requestUrl:    "/get_funky_content_length_body_hello",
		contentLength: unused,
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
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_POST,
		requestUrl:    "/post_identity_body_world?q=search#hey",
		contentLength: unused,
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
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_POST,
		requestUrl:    "/post_chunked_all_your_base",
		contentLength: unused,
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
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_POST,
		requestUrl:    "/two_chunks_mult_zero_end",
		contentLength: unused,
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
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_POST,
		requestUrl:    "/chunked_w_trailing_headers",
		contentLength: unused,
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
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_POST,
		requestUrl:    "/chunked_w_nonsense_after_length",
		contentLength: unused,
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
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_POST,
		requestUrl:    "/with_\"stupid\"_quotes?foo=\"bar\"",
		contentLength: unused,
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
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            0,
		//method: HTTP_POST,
		requestUrl:    "/test",
		contentLength: unused,
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
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_POST,
		requestUrl:    "/test.cgi?foo=bar?baz",
		contentLength: unused,
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
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_POST,
		requestUrl:    "/test",
		contentLength: unused,
	},
	/*
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
			messageCompleteOnEof: false,
			httpMajor:            1,
			httpMinor:            1,
			//method: HTTP_POST,
			requestUrl:    "/demo",
			contentLength: unused,
			headers: [][2]string{{"Host", "example.com"},
				{"Connection", "Upgrade"},
				{"Sec-WebSocket-Key2", "12998 5 Y3 1  .P00"},
				{"Sec-WebSocket-Protocol", "sample"},
				{"Upgrade", "WebSocket"},
				{"Sec-WebSocket-Key1", "4 @1  46546xW%0l 1 5"},
				{"Origin", "http://example.com"},
			},
		},
	*/
}

var responses = []message{
	{
		name:  "HTTP/1.1 with chunked endocing and a 200 response",
		hType: RESPONSE,
		raw: "HTTP/1.1 200 OK\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"\r\n" +
			"0\r\n" +
			"\r\n",
		statusCode:              200,
		responseStatus:          "OK",
		shouldKeepAlive:         true,
		messageCompleteOnEof:    false,
		messageCompleteCbCalled: true,
		httpMajor:               1,
		httpMinor:               1,
		//method: HTTP_GET,
		contentLength: unused,
		headers: [][2]string{
			{"Transfer-Encoding", "chunked"},
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
		messageCompleteOnEof:    false,
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
		messageCompleteOnEof:    false,
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
		messageCompleteOnEof:    false,
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
		messageCompleteOnEof:    false,
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
		messageCompleteOnEof:    false,
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
		messageCompleteOnEof:    false,
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
		messageCompleteOnEof:    false,
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
		messageCompleteOnEof:    false,
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
		messageCompleteOnEof:    false,
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
		messageCompleteOnEof:    false,
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
	MessageBegin: func(p *Parser) {
		m := p.GetUserData().(*message)
		m.messageBeginCbCalled = true
		m.hType = p.hType
	},
	URL: func(p *Parser, url []byte) {
		m := p.GetUserData().(*message)
		m.requestUrl += string(url)
	},
	Status: func(p *Parser, status []byte) {
		m := p.GetUserData().(*message)
		m.responseStatus += string(status)
	},
	HeaderField: func(p *Parser, headerField []byte) {
		m := p.GetUserData().(*message)
		m.headers = append(m.headers, [2]string{string(headerField), ""})
	},
	HeaderValue: func(p *Parser, headerValue []byte) {
		m := p.GetUserData().(*message)
		m.headers[len(m.headers)-1][1] = string(headerValue)
	},
	HeadersComplete: func(p *Parser) {
		m := p.GetUserData().(*message)
		m.headersCompleteCbCalled = true
	},
	Body: func(p *Parser, body []byte) {
		m := p.GetUserData().(*message)
		m.body += string(body)
	},
	MessageComplete: func(p *Parser) {
		m := p.GetUserData().(*message)
		m.messageCompleteCbCalled = true
		m.httpMajor = uint16(p.Major)
		m.httpMinor = uint16(p.Minor)

	},
}

func parse(p *Parser, data string) (int, error) {
	return p.Execute(&settingTest, []byte(data))
}

func test_Message(t *testing.T, m *message) {
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

		if msg1len > 0 {
			n1, err1 = parse(p, msg1Message)
			assert.NoError(t, err1)
			// 如果有upgrade状态, 就不需要再重复送往数据
			if got.messageCompleteCbCalled && p.Upgrade {
				got.upgrade = msg1Message[n1:]
				msg1Message = ""
			} else {

				msg1Message = msg1Message[n1:]
			}
		}

		data = msg1Message + msg2Message
		n2, err = parse(p, data)
		if got.messageCompleteCbCalled && p.Upgrade {
			got.upgrade += data[n2:]
			goto test
		}
		assert.NoError(t, err)

	test:

		_, err = parse(p, "")
		assert.NoError(t, err)
		if b := m.eq(t, got); !b {
			t.Logf("msg1.len:%d, msg2.len:%d, test case name:%s\n", len(msg1Message), len(msg2Message), m.name)
			t.Logf("msg1len:%d, msg1(%s)", msg1len, msg1Message)
			t.Logf("msg2(%s)", msg2Message)
			t.Logf("upgrade:%t, got.messageCompleteCbCalled:%t, data:(%s)", p.Upgrade, got.messageCompleteCbCalled, data[n2:])
			break
		}

	}
}

func Test_Message(t *testing.T) {
	for _, req := range requests {
		test_Message(t, &req)
		_ = req
	}

	for _, rsp := range responses {
		test_Message(t, &rsp)
		_ = rsp
	}
}

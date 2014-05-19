package server_test

import (
	//"github.com/vizidrix/gocqrs/net/server"
	"testing"
)

func Test_Should_return_views_list(t *testing.T) {
	request :=
`GET /api/v1/views HTTP/1.1
Host: localhost:8080
Connection: keep-alive
Cache-Control: max-age=0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.137 Safari/537.36
Accept-Encoding: gzip,deflate,sdch
Accept-Language: en-US,en;q=0.8`

	
	t.Errorf("Should serve views [ \n%s\n ]\n", request)
}

func Test_Should_not_serve_file_requests(t *testing.T) {
	request :=
`GET /favicon.ico HTTP/1.1
Host: localhost:8080
Connection: keep-alive
Accept: */*
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.137 Safari/537.36
Accept-Encoding: gzip,deflate,sdch
Accept-Language: en-US,en;q=0.8`

	t.Errorf("Should not serve [ \n%s\n ]\n", request)
}

func Test_Should_accept_post_data(t *testing.T) {
	request :=
`POST /api/v1/things/addcommand HTTP/1.1
Host: localhost:8080
Connection: keep-alive
Content-Length: 14
Cache-Control: max-age=0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Origin: http://localhost:8080
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.137 Safari/537.36
Content-Type: application/x-www-form-urlencoded
Referer: http://localhost:8080/
Accept-Encoding: gzip,deflate,sdch
Accept-Language: en-US,en;q=0.8

blabha=adsfasd`

	t.Errorf("Should read post data [ \n%s\n \n", request)
}
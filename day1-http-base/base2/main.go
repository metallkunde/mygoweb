package main

// $ curl http://localhost:9999/
// URL.Path = "/"
// $ curl http://localhost:9999/hello
// Header["Accept"] = ["*/*"]
// Header["User-Agent"] = ["curl/7.54.0"]
// curl http://localhost:9999/world
// 404 NOT FOUND: /world

import (
	"fmt"
	"log"
	"net/http"
)

// 定义了一个空的结构体Engine，实现了方法ServeHTTP。
// 这个方法有2个参数，第二个参数是 Request ，该对象包含了该HTTP请求的所有的信息，
// 比如请求地址、Header和Body等信息；第一个参数是 ResponseWriter ，
// 利用 ResponseWriter 可以构造针对该请求的响应。
type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":19999", engine))
}

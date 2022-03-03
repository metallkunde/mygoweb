package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// context必要性

/* 	对Web服务来说，无非是根据请求*http.Request，构造响应http.ResponseWriter。
 	但是这两个对象提供的接口粒度太细，比如我们要构造一个完整的响应，需要考虑消息头(Header)
	和消息体(Body)，而 Header 包含了状态码(StatusCode)，消息类型(ContentType)等几乎
	每次请求都需要设置的信息。因此，如果不进行有效的封装，那么框架的用户将需要写大量重复，
	繁杂的代码，而且容易出错。针对常用场景，能够高效地构造出 HTTP 响应是一个好的框架必须考虑的点。

	Context随着每一个请求的出现而产生，结束而销毁，和当前请求强相关的信息都应由Context承载。
 	因此，设计 Context 结构，扩展性和复杂性留在了内部，而对外简化了接口。路由的处理函数，
	以及将要实现的中间件，参数都统一使用 Context 实例.
*/

// H 为map[string]interface{}别名，构建JSON数据时，显得更简洁。
type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求信息直接访问
	Path   string
	Method string
	Params map[string]string
	// 响应信息
	StatusCode int
	// 中间件
	handlers []HandlerFunc
	index    int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// PostForm 提供了访问PostForm参数的方法
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 提供了访问Query参数的方法。
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

//提供了快速构造String/Data/JSON/HTML响应的方法。

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

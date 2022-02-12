package gee

//将router相关的代码独立后，gee.go简单了不少。
//最重要的还是通过实现了 ServeHTTP 接口，接管了所有的 HTTP 请求。

//将和路由相关的方法和结构提取了出来，放到了一个新的文件中router.go，
//方便我们下一次对 router 的功能进行增强，例如提供动态路由的支持。
//router 的 handle 方法作了一个细微的调整，即 handler 的参数，变成了 Context。
import (
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

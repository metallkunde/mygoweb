package gee

import "net/http"

//将和路由相关的方法和结构提取了出来，方便我们下一次对 router 的功能进行增强，
//例如提供动态路由的支持。 router 的 handle 方法作了一个细微的调整，
//即 handler 的参数，变成了 Context。

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// 路由映射表中按照其key-value注册方法
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

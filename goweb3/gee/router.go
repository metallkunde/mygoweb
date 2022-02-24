package gee

import (
	"net/http"
	"strings"
)

/*
	将和路由相关的方法和结构提取了出来，方便我们下一次对 router 的功能进行增强，例如提供动态路由的支持。
	router 的 handle 方法作了一个细微的调整，即 handler 的参数，变成了 Context。
	使用 roots 来存储每种请求方式的Trie 树根节点。
	使用 handlers 存储每种请求方式的 HandlerFunc 。
*/
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 去除路由中的斜杠，并将每层信息都存入slice，遇到*就停止
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	//取消首空格，并去除*后所有内容
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 路由映射表中按照其key-value注册方法
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

/*
	getRoute 函数中，还解析了:和*两种匹配符的参数，返回一个 map 。
	例如/p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}，
	/static/css/geektutu.css匹配到/static/*filepath，解析结果为
	{filepath: "css/geektutu.css"}。
*/
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

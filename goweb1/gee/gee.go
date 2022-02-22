package gee

import (
	"fmt"
	"log"
	"net/http"
)

//在实现Engine之前，我们调用 http.HandleFunc 实现了路由和Handler的映射，
//也就是只能针对具体的路由写处理逻辑

// HandlerFunc 定义路由映射的处理方法
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine 中定义路由映射表router
// 其中key为请求方法+静态路由地址，value为用户映射的处理方法
// 拦截了所有的HTTP请求，拥有了统一的控制入口。
// 在这里我们可以自由定义路由映射的规则，
// 也可以统一添加一些处理逻辑，例如日志、异常处理等。
type Engine struct {
	router map[string]HandlerFunc
}

// New 初始化结构体Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 路由映射表中按照其key-value注册方法
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	log.Printf("Route %4s - %s", method, pattern)
	engine.router[key] = handler
}

// GET 将路由和处理方法注册到映射表 router
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run 包装ListenAndServe
// 第一个参数是地址，第二个参数则代表处理所有的HTTP请求的实例，
// 若为nil则代表使用标准库中的实例处理。
// 第二个参数是基于net/http标准库实现Web框架的入口。
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP 方法的作用为解析请求的路径，查找路由映射表，
// 如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, r)
	} else {
		fmt.Fprintf(w, "404 not found:%s\n", r.URL)
	}
}

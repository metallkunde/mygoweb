package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

/*	在实现Engine之前，我们调用 http.HandleFunc 实现了路由和Handler的映射，
	也就是只能针对具体的路由写处理逻辑
*/

// HandlerFunc 定义路由映射的处理方法
type HandlerFunc func(c *Context)

type RouterGroup struct {
	//前缀
	prefix string
	//中间件是应用在分组上的，存储应用在该分组上的中间件
	middlewares []HandlerFunc
	//分组嵌套
	parent *RouterGroup
	/*
		Group对象，还需要有访问Router的能力，为了方便可以在Group中保存一个指针指向Engine，
		整个框架的所有资源都是由Engine统一协调的，那么就可以通过Engine间接地访问各种接口了。
	*/
	engine *Engine
}

// Engine 中定义路由映射表router
// 其中key为请求方法+静态路由地址，value为用户映射的处理方法
// 拦截了所有的HTTP请求，拥有了统一的控制入口。
// 在这里我们可以自由定义路由映射的规则，
// 也可以统一添加一些处理逻辑，例如日志、异常处理等。
// 将Engine作为最顶层的分组，Engine拥有RouterGroup所有的能力。
type Engine struct {
	*RouterGroup
	router *router
	//存储所有的group
	groups []*RouterGroup
	//将所有的模板加载进内存
	htmlTemplates *template.Template
	//所有的自定义模板渲染函数
	funcMap template.FuncMap
}

// New 初始化结构体Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Default 默认实例使用 Logger 和 Recovery 中间件
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

// Group 创建新的RouterGroup，所以groups拥有同一个instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Use 向分组添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// 路由映射表中按照其key-value注册方法
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET 将路由和处理方法注册到映射表 router
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static 将磁盘上的某个文件夹root映射到路由relativePath
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

// for custom render function
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// Run 包装ListenAndServe
// 第一个参数是地址，第二个参数则代表处理所有的HTTP请求的实例，若为nil则代表使用标准库中的实例处理。
// 第二个参数是基于net/http标准库实现Web框架的入口。
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP 方法的作用为解析请求的路径，查找路由映射表，
// 如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND
// 当我们接收到一个具体请求时，要判断该请求适用于哪些中间件，
// 在这里我们简单通过 URL 的前缀来判断。得到中间件列表后，赋值给 c.handlers
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, r)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

package gac

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

//gac处理路由的func别名

type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup                    //让Engine具有RouterGroup的功能
	router       *router            //路由信息
	groups       []*RouterGroup     //所有的分组
	htmlTemplate *template.Template //处理html渲染
	funcMap      template.FuncMap   //处理html渲染
}

//分组控制
type RouterGroup struct {
	prefix     string        //分组前缀
	middleware []HandlerFunc //中间件
	engine     *Engine
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middleware []HandlerFunc
	for _, group := range engine.groups {
		//添加用户自定义的中间件handler
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middleware = append(middleware, group.middleware...)
		}
	}
	//实例化Context
	c := newContext(w, req)
	c.handlers = middleware //还没加路由handler
	c.engine = engine
	engine.router.handle(c)
}

func New() *Engine {
	//初始一个空路由信息
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRouter(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middleware ...HandlerFunc) {
	group.middleware = append(group.middleware, middleware...)
}

//静态资源Handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	// 用‘/’连接两部分字符
	absolutePath := path.Join(group.prefix, relativePath)
	//相当于加了一个前缀
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

//处理静态资源,将本地磁盘中的文件夹root,映射到relativePath相对路径
//eg g.Static("/static","C:/User/ddd/resources") 这样resources下的文件就可以通过
// ip:port/static/js/login.js 访问 C:/User/ddd/resources/js/login.js
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	//相当于变成了 /static/*filepath 也就是我们前面写过的动态路由
	urlPattern := path.Join(relativePath, "*filepath")
	//正常的GET请求
	group.GET(urlPattern, handler)
}

//自定义渲染函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

//加载模板
func (engine *Engine) LoadHtmlGlob(pattern string) {
	engine.htmlTemplate = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

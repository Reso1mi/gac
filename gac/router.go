package gac

import (
	"log"
	"net/http"
	"strings"
)

// roots key eg, roots['GET'] roots['POST'] 用method做Trie的根节点
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

//构建路由信息
func newRouter() *router {
	return &router{
		handlers: map[string]HandlerFunc{},
		roots:    map[string]*node{},
	}
}

//将整体的路由信息以’/‘划分为part []string
func parsePattern(pattern string) []string {
	//和Java不一样,go的Split并不会过滤尾部的空字符
	ps := strings.Split(pattern, "/")
	part := make([]string, 0)
	for _, p := range ps {
		if p != "" { //过滤空字符
			part = append(part, p)
			if p[0] == '*' { //遇到*就结束,只允许一个*
				break
			}
		}
	}
	return part
}

func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route: %4s - %s", method, pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		//method作为根节点
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parsePattern(pattern), 0)
	r.handlers[key] = handler
}

//获取路由在Trie中完整信息的同时,解析URL中的动态路由参数
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	//将URL解析成路由part做搜索
	searchPart := parsePattern(path)
	//构建动态路由对应的param
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	//在Trie中搜索路由信息
	n := root.search(searchPart, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for i, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchPart[i]
			}
			if part[0] == '*' && len(part) > 1 {
				//用 ’/‘ 连接 searchPart[i:]中的所有字符
				params[part[1:]] = strings.Join(searchPart[i:], "/")
				//遇到一个*就结束, Only one * is Allowed
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if node != nil {
		c.Params = params
		key := c.Method + "-" + node.pattern
		//r.handlers[key](c) //用户的Handler处理请求
		//将路由Handler加入整个handlers中
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		//c.String(http.StatusNotFound, "404 NOT FOUNT: %s\n Power By Gac !", c.Path)
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUNT: %s\n Power By Gac !", c.Path)
		})
	}
	c.Next() //将执行权交给下一个Handler
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	node := make([]*node, 0)
	root.travel(&node)
	return node
}

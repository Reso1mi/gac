package gac

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//JSON串封装
type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	//请求路径
	Path string
	//请求方式GET POST
	Method string
	//请求的动态路由参数
	Params map[string]string
	//状态码
	StatusCode int
	//用户自定义的middleware和路由handler
	handlers []HandlerFunc
	index    int
	engine   *Engine
}

//封装一下
func (c *Context) Param(key string) string {
	return c.Params[key]
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	//c.handlers[c.index](c)
	//并不是所有的handler都调用了Next(),所以这里需要遍历一下所有的handlers去依次调用
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) PostForm(key string) string {
	//URL和Form组合后的第一个值,貌似是Form优先
	//对比PostFormValue只获取Form的值
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	//写入状态码
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

//String/JSON/Data/HTML的响应
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

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	//c.Writer.Write([]byte(html))
	if err := c.engine.htmlTemplate.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(http.StatusInternalServerError, err.Error())
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

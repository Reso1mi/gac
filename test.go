package main

import (
	"fmt"
	"gac"
	"html/template"
	"net/http"
	"strings"
	"time"
)

//错误恢复
func main() {
	g := gac.Default()
	g.GET("/", func(c *gac.Context) {
		c.String(http.StatusOK, "Hello Resolmi!!!")
	})
	g.GET("/painc", func(c *gac.Context) {
		names := []string{"imlgw"}
		//越界
		c.String(http.StatusOK, names[10000])
	})
	g.Run(":8888")
}

type student struct {
	Name string
	Age  int8
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

//HTML模板
func main3() {
	r := gac.New()
	r.Use(gac.Logger())
	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	r.LoadHtmlGlob("template/*")
	//将相对路径 ./static 映射到 /asserts
	r.Static("/assets", "./static")
	stu1 := &student{Name: "Resolmi", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *gac.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *gac.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gac.H{
			"title":  "gac",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *gac.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gac.H{
			"title": "gee",
			"now":   time.Date(2020, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}

//分组路由
func main2() {
	r := gac.New()
	r.Use(gac.Logger())
	v1 := r.Group("/v1")
	{

		v1.GET("/hello", func(c *gac.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	v2.Use(func(c *gac.Context) {
		fmt.Println("Only For V2")
	})
	{
		v2.GET("/hello/:name", func(c *gac.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gac.Context) {
			c.JSON(http.StatusOK, gac.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}
	r.Run(":8888")
}

func myTTTT() {
	//var arr [2]int
	arr := make([]int, 2)
	test(arr)
	fmt.Println(arr)

	var arr2 [2]int
	test2(&arr2)
	fmt.Println(arr2)

	//var su = new(stu)
	var su = &stu{
		name: "TEST",
	}
	test3(su)
	fmt.Println(su)

	//ss := strings.Split(",aa,bb,cc,dd,,,", ",") //8
	ss := strings.Split("/", "/") //2
	for i := range ss {
		fmt.Println(ss[i])
	}
	fmt.Println(len(ss))

	//fmt.Println(strings.Join([]string{"dasdas", "css", "ac"}, "-")) //dasdas-css-ac
	var aa rune = 'a'
	var bb rune = 'b'

	var sa string = "a"
	var sb string = "b"

	fmt.Println(aa - bb)
	fmt.Println(int8(sa[0]) - int8(sb[0]))
}

type stu struct {
	name string
}

func test3(su *stu) {
	su.name = "IMLGW.TOP"
}

func test2(arr *[2]int) {
	arr[0] = 10010101001
}

func test(arr []int) {
	arr[0] = 10010101001

}

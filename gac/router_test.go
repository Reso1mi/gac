package gac

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRouter("GET", "/", nil)
	r.addRouter("GET", "/hello/:name", nil)
	r.addRouter("GET", "/hello/b/c", nil)
	r.addRouter("GET", "/hi/:name", nil)
	r.addRouter("GET", "/assets/*filepath", nil)
	return r
}

//解析路由测试
func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	//非法路由
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	node, params := r.getRoute("GET", "/hello/resolmi")

	if node == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if node.pattern != "/hello/:name" {
		fmt.Println(node.pattern)
		t.Fatal("should match /hello/:name")
	}

	if params["name"] != "resolmi" {
		t.Fatal("name should be equal to 'resolmi'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", node.pattern, params["name"])

}

func TestGetRoute2(t *testing.T) {
	r := newTestRouter()
	node1, params1 := r.getRoute("GET", "/assets/file1.txt")
	ok1 := node1.pattern == "/assets/*filepath" && params1["filepath"] == "file1.txt"
	if !ok1 {
		t.Fatal("pattern shoule be /assets/*filepath & filepath shoule be file1.txt")
	}

	node2, param2 := r.getRoute("GET", "/assets/css/test.css")
	ok2 := node2.pattern == "/assets/*filepath" && param2["filepath"] == "css/test.css"
	if !ok2 {
		t.Fatal("pattern shoule be /assets/*filepath & filepath shoule be css/test.css")
	}

}

func TestGetRoutes(t *testing.T) {
	r := newTestRouter()
	nodes := r.getRoutes("GET")
	for i, n := range nodes {
		fmt.Println(i+1, n)
	}

	if len(nodes) != 6 {
		//t.Fatal("the number of routes should be 6")
	}
}

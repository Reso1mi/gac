package gac

import "strings"

/*
	和一般的字典树还不太一样,有一点改进
*/
type node struct {
	pattern   string  //路由信息
	part      string  //路由信息的各个部分
	children  []*node //子节点
	isDynamic bool    //是否是动态路由，包含*和：
}

//找到一个以n为根节点和part相等的子节点(只需要一个,用于插入的)
// ip/a/* == ip/a/b/c
func (n *node) match(part string) *node {
	for _, child := range n.children {
		if child.isDynamic || child.part == part {
			return child
		}
	}
	return nil
}

//找到所有以n为根节点,和part相等的子节点
func (n *node) matchAll(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.isDynamic || child.part == part {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//插入新路由信息
func (n *node) insert(pattern string, parts []string, index int) {
	if index == len(parts) {
		//给最后一层的pattern赋值
		n.pattern = pattern
		return
	}
	cur := parts[index]
	//匹配的孩子节点
	child := n.match(cur)
	if child == nil { //没有匹配的了，说明是新节点，匹配结束，将当前part加入当前层的child
		child = &node{part: cur, isDynamic: cur[0] == ':' || cur[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, index+1)
}

//搜索路由信息,如果存在路由信息就返回该路由在Trie中的最后一个节点
func (n *node) search(parts []string, index int) *node {
	if len(parts) == index || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	cur := parts[index]
	children := n.matchAll(cur)

	for _, child := range children {
		result := child.search(parts, index+1)
		if result != nil {
			return result
		}
	}

	return nil
}

func (n *node) travel(nodes *[]*node) {
	if n != nil && n.pattern != "" {
		*nodes = append(*nodes, n)
	}
	children := n.children
	for _, child := range children {
		child.travel(nodes)
	}
}

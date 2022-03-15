package gee

import (
	"fmt"
	"strings"
)

type node struct {
	//待匹配路由，例如 /p/:lang
	pattern string
	//路由中的一部分，例如 :lang
	part string
	//子节点，例如 [doc, tutorial, intro]
	children []*node
	//是否精确匹配，part 含有 : 或 * 时为true
	isWild bool
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

/*
	递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个，有一点需要注意，
	/p/:lang/doc只有在第三层节点，即doc节点，pattern才会设置为/p/:lang/doc。
	p和:lang节点的pattern属性皆为空。因此，当匹配结束时，我们可以使用
	n.pattern == ""来判断路由规则是否匹配成功。
*/
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

/*
	递归查询每一层的节点，退出规则是，匹配到了*，匹配失败，或者匹配到了第len(parts)层节点。
*/
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

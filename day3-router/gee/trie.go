package gee

import (
	"strings"
)

type node struct {
	pattern  string  //待匹配路由，例如/p/:lang
	part     string  //路由中的一部分，例如 :lang
	children []*node //子节点，例如[doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part含有：或者*为true
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			//如果遇见当前为*的话直接break
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

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
			//如果第一个字符为*，并且不仅仅是*， 那么记录后面的所有参数即可
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 第一个匹配成功的节点，用于插入
func (n *node) mathChild(part string) *node {
	for _, child := range n.children {
		//如果直接等于或者当前点是模糊匹配的，直接返回child即可
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	//感觉不必如此，可以直接等到最后一个遍历完，然后设置一个标记即可
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.mathChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

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
		res := child.search(parts, height+1)
		if res != nil {
			return res
		}
	}
	return nil
}

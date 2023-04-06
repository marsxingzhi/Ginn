package ginn

import "strings"

/**
前缀树的实现

1. 节点的定义
2. 插入
3. 查询

												/
					/:lang						/about					/ping
			/intro	/tuorial	/doc							/blog			/related
*/

type TrieNode struct {
	path      string      // 从根节点到当前节点的字符串路径，例如：/v1/ping
	val       string      // 当前节点的元素值，例如：/ping
	children  []*TrieNode // 孩子节点，例如：[/intor, /tuorial, /doc]
	isDynamic bool        // 是否是动态路由，例如：包含:或者*的，如/:lang
}

// 匹配child
func (node *TrieNode) matchChild(val string) (*TrieNode, bool) {
	for _, child := range node.children {
		if child.isDynamic || child.val == val {
			return child, true
		}
	}
	return nil, false
}

// 匹配children
func (node *TrieNode) matchChildren(val string) ([]*TrieNode, bool) {
	res := make([]*TrieNode, 0)

	for _, child := range node.children {
		if child.isDynamic || child.val == val {
			res = append(res, child)
		}
	}

	return res, len(res) > 0
}

/**
pattenr：待添加的路由
vals：将待添加的路由分割成多个部分
index：当前处理的高度

思路：
1.
*/
func (node *TrieNode) insert(path string, vals []string, index int) {
	// 递归结束
	if len(vals) == index {
		// 更新当前节点的路径
		node.path = path
		return
	}
	val := vals[index]
	child, ok := node.matchChild(val)
	if !ok {
		// 未匹配到，则创建一个节点
		child = &TrieNode{
			val:       val,
			isDynamic: isDynamic(val),
		}
		node.children = append(node.children, child)
	}
	// 再次执行inset方法，其实就是为了更新最终的pattern
	// 假设现在是叶子节点，那么再次调用insert方法，会在下一次函数开始的if判断中，判断成功，并直接返回
	child.insert(path, vals, index+1)
}

func (node *TrieNode) search(vals []string, index int) (*TrieNode, bool) {
	if len(vals) == index || strings.HasPrefix(node.val, "*") {
		if node.path == "" {
			return nil, false
		}
		return node, true
	}

	val := vals[index]
	children, ok := node.matchChildren(val)
	if !ok {
		return nil, false
	}
	for _, child := range children {
		res, ok := child.search(vals, index+1)
		if ok {
			return res, ok
		}
	}
	return nil, false
}

func isDynamic(val string) bool {
	return val[0] == ':' || val[0] == '*'
}

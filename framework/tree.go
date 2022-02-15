package framework

import (
	"errors"
	"log"
	"strings"
)

type Tree struct {
	root *node // 根结点
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{root}
}

type node struct {
	isLast   bool                // 表明该结点是否为叶子节点
	segment  string              // 当前结点的
	handlers []ControllerHandler // 中间件+控制器
	childs   []*node             // 当前结点的子节点
	parent   *node               // 父节点，双向指针
}

func newNode() *node {
	return &node{
		isLast:  false,
		segment: "",
		childs:  []*node{},
		parent:  nil,
	}
}

func (tree *Tree) AddRouter(uri string, handlers []ControllerHandler) error {
	n := tree.root
	// 完全匹配，说明冲突
	if n.matchNode(uri) != nil {
		return errors.New("route exist: " + uri)
	}
	segments := strings.Split(uri, "/")
	for index, segment := range segments {
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		isLast := index == len(segments)-1
		var objNode *node
		childNodes := n.filterChildNodes(segment)
		// 如果有匹配的子节点
		if len(childNodes) > 0 {
			// 如果有segment相同的子节点，则选择这个子节点
			for _, cnode := range childNodes {
				if cnode.segment == segment {
					objNode = cnode
					break
				}
			}
		}
		if objNode == nil {
			cnode := newNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handlers = handlers
			}
			// 父节点指针修改
			cnode.parent = n
			n.childs = append(n.childs, cnode)
			objNode = cnode
		}
		n = objNode
	}
	return nil
}

// 匹配uri
func (tree *Tree) FindHandler(uri string) []ControllerHandler {
	matchNode := tree.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}
	return matchNode.handlers
}

// 判断一个segment是否是通用segment，即以:开头
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

// 过滤下一层满足segment规则的子节点
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}
	// 如果segment是通配符，则所有下一层子节点都满足需求
	if isWildSegment(segment) {
		return n.childs
	}
	// 过滤所有的下一层子节点
	nodes := make([]*node, 0, len(n.childs))
	for _, cnode := range n.childs {
		if isWildSegment(cnode.segment) {
			nodes = append(nodes, cnode) // 如果下一层子节点有通配符，则满足需求
		} else if cnode.segment == segment {
			nodes = append(nodes, cnode) // 如果下一层子节点没有通配符，但是文本完全匹配，则满足需求
		}
	}
	return nodes
}

// 判断路由是否已经在节点的所有子节点树中存在了   => 用于建立or查找路由表？
func (n *node) matchNode(uri string) *node {
	// 使用分隔符将uri切割为两个部分
	segments := strings.SplitN(uri, "/", 2)
	// 第一个部分用于匹配下一层子节点
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}
	// 匹配符合的下一层子节点
	cnodes := n.filterChildNodes(segment)

	log.Println("segment", segment)
	for _, n := range cnodes {
		log.Print(n.segment, " ")
	}
	log.Println()

	if cnodes == nil || len(cnodes) == 0 {
		return nil
	}

	// 如果只有一个segment，则是最后一个标记
	if len(segments) == 1 {
		// 如果segment已经是最后一个节点，判断这些cnode是否有isLast标志
		for _, tn := range cnodes {
			if tn.isLast {
				return tn
			}
		}
		// 都不是最后一个节点->nil
		return nil
	}

	// 如果有2个segment，递归每个子节点继续查找
	for _, tn := range cnodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}
	return nil
}

// 将uri解析为params -> 查找出整个匹配链路中的通配符节点和对应 URI 中的分段
func (n *node) parseParamsFromEndNode(uri string) map[string]string {
	ret := map[string]string{}
	segments := strings.Split(uri, "/")
	cnt := len(segments)
	cur := n
	for i := cnt - 1; i >= 0; i-- {
		if cur.segment == "" {
			break
		}
		if isWildSegment(cur.segment) {
			// 设置params
			ret[cur.segment[1:]] = segments[i]
		}
	}
	return ret
}

package go_radix

import (
	. "github.com/fuyao-w/common-util"
	"sort"
	"strings"
)

/*
radix 树实现

TODO 把 edge.parent 干掉减小复杂度
*/
type (
	// node 树的节点
	node struct {
		predecessor *edge               // 节点的前继边
		successor   []*edge             // 节点的后继边，按 edge.prefix[0] 从小到大排序
		leaf        *Tuple[string, any] // 当前节点是否有元素
	}
	// edge 树的边
	edge struct {
		parent, child *node
		prefix        string // 每个边保存局部的前缀，从顶点到某一叶子节点的局部前缀组成完整的 key
	}
	// Tree 通过树对外提供 API
	Tree struct {
		root  *node
		count uint64
	}
)

// New 创建一个新的 radix 树
func New() *Tree {
	return &Tree{root: new(node)}
}

func (e *edge) hasLeaf() bool {
	return e.child.leaf != nil
}

func (e *edge) setLeaf(key string, value any) {
	e.child.setLeaf(key, value)
}
func (e *edge) delLeaf() {
	e.child.leaf = nil
}

func (e *edge) getLeaf() *Tuple[string, any] {
	return e.child.leaf
}
func (n *node) setLeaf(key string, value any) {
	n.leaf = Ptr(BuildTuple(key, value))
}

func (n *node) isLeaf() bool {
	return n.leaf != nil
}

// longestPrefix 比较两个字符串的公共前缀索引
func longestPrefix(item, target string) (idx int) {
	for i := range item {
		if !strings.HasPrefix(target, item[:i+1]) {
			return i
		}
	}
	return len(item)
}

// findEdge 寻找某一条边如果不存在则 found 返回 false，idx 是 target 的索引位置
func findEdge(s []*edge, target string) (idx int, found bool) {
	return sort.Find(len(s), func(i int) int {
		return -strings.Compare(s[i].prefix[:1], target[0:1])
	})
}

// appendEdge 追加一条边，按照 edge.prefix[0] 从小到大的位置进行插入
func (n *node) appendEdge(e *edge) {
	e.parent = n
	idx, _ := findEdge(n.successor, e.prefix)
	n.successor = append(n.successor[:idx], append([]*edge{e}, n.successor[idx:]...)...)
}

// Insert 插入元素 ok 为 false 代表跟新 ，old 返回之前的 value
func (t *Tree) Insert(key string, value any) (old any, ok bool) {
	return t.insert(key, value)
}

func (t *Tree) insert(key string, value any) (old any, ok bool) {
	curNode := t.root
	search := key

	for {
		// 如果 prefix 是 search 的完全子串，直接终止
		if len(search) == 0 {
			if !curNode.isLeaf() {
				t.count++
			} else {
				ok = true
				old = curNode.leaf.B
			}
			curNode.setLeaf(key, value)
			return
		}
		// 如果 prefix 是 search 的完全子串，直接终止
		idx, found := findEdge(curNode.successor, search)
		if !found {
			next := newEdge(search, nil)
			next.setLeaf(key, value)
			curNode.appendEdge(next)
			t.count++
			return
		}
		curEdge := curNode.successor[idx]
		matchIndex := longestPrefix(curEdge.prefix, search)
		if matchIndex == 0 {
			panic("edge not have common prefix")
		}

		oriPrefix, newPrefix := search[:matchIndex], search[matchIndex:]

		// 完全匹配,直接返回
		//if len(oriPrefix) == len(curEdge.prefix) && len(newPrefix) == 0 {
		//	if !curEdge.hasLeaf() {
		//		t.count++
		//	}
		//	curEdge.setLeaf(key, value)
		//	return
		//}
		search = newPrefix
		//if oriPrefix == curEdge.prefix && len(search) > 0 {
		//	// prefix 是 search 的完全子串，直接将 search[matchIndex:] 插入子节点
		//	curNode = curEdge.child
		//} else
		if oriPrefix != curEdge.prefix { // prefix 不是 search 子串 (部分重叠)，先拆分，再确定是否增加子节点
			//创建新的边，并且将当前边的 leaf 移动到新边
			next := newEdge(curEdge.prefix[len(oriPrefix):], curEdge.getLeaf())
			next.child.successor = append([]*edge(nil), curEdge.child.successor...)
			next.parent = curEdge.child
			// 将当前边的下游节点后继更新成 next , leaf 置空 ，更新 prefix
			//curEdge.child = next.parent
			curEdge.child.successor = append([]*edge(nil), next)
			curEdge.delLeaf()
			curEdge.prefix = oriPrefix
		}
		if len(search) > 0 {
			curNode = curEdge.child
		}
	}
}

// newEdge 创建一条新的边
func newEdge(localPrefix string, leaf *Tuple[string, any]) *edge {
	edge := &edge{prefix: localPrefix}
	downNode := &node{
		predecessor: edge,
		leaf:        leaf,
	}
	edge.child = downNode
	return edge
}

// walk 广度遍历节点 n，walkFn 返回 true 的时候会继续遍历，否则会停止
func (t *Tree) walk(n *node, walkFn func(key string, val any) bool) {
	if n == nil {
		return
	}
	var (
		successor, nextSuccessor []*edge
	)
	if n.predecessor == nil && n.isLeaf() { // 根节点
		successor = append(successor, newEdge("", n.leaf))
	}
	successor = append(successor, n.successor...)
	for len(successor) > 0 {
		nextSuccessor = nil
		for _, e := range successor {
			if e.hasLeaf() {
				leaf := e.getLeaf()
				if !walkFn(leaf.A, leaf.B) {
					return
				}
			}
			nextSuccessor = append(nextSuccessor, e.child.successor...)
		}
		successor = nextSuccessor
	}
}

// Walk 从匹配 prefix 最长前缀的节点开始遍历树，walkFn 返回 true 代表继续遍历，否则终止
func (t *Tree) Walk(prefix string, walkFn func(key string, val any) bool) {
	t.walk(t.findLongestPrefixNode(prefix), walkFn)
}

// findLongestPrefixNode 找到具有最长公共前缀的节点
func (t *Tree) findLongestPrefixNode(prefix string) *node {
	curNode := t.root
	search := prefix
	if len(prefix) == 0 {
		return curNode
	}
	var successor = curNode.successor
	for len(successor) > 0 {
		idx, found := findEdge(successor, search)
		if !found {
			break
		}
		curEdge := curNode.successor[idx]
		curNode = curEdge.child
		if len(curEdge.prefix) >= len(search) {
			return curNode
		}
		successor = curNode.successor
		search = search[len(curEdge.prefix):]
	}
	return nil
}

// Get 等值查找
func (t *Tree) Get(key string) (val any, exist bool) {
	node := t.findLongestPrefixNode(key)
	if node.isLeaf() {
		return node.leaf.B, true
	}
	return
}

// Len 树的有效元素数量
func (t *Tree) Len() uint64 {
	return t.count
}

// LongestPrefix 返回有效的最长公共前缀
func (t *Tree) LongestPrefix(prefix string) (longestPrefix string, value interface{}, exist bool) {
	node := t.findLongestPrefixNode(prefix)
	if node == nil {
		return
	}
	if node.isLeaf() {
		return node.leaf.A, node.leaf.A, true
	}
	return
}

// Delete 等值删除
func (t *Tree) Delete(key string) (ok bool) {
	n := t.findLongestPrefixNode(key)
	if n == nil {
		return false
	}
	if !n.isLeaf() || n.leaf.A != key {
		return false
	}
	n.leaf = nil
	if len(n.successor) == 0 { // 如果是叶子边，可以直接删除
		n.deleteEdge(n.predecessor)
	}
	if n.predecessor != nil {
		n.predecessor.parent.mergeChild()
	}
	t.count--
	return ok
}

// deleteEdge 将当前节点的前继边从父节点上删除
func (n *node) deleteEdge(target *edge) {
	if target == nil {
		return
	}
	successor := n.predecessor.parent.successor
	if len(successor) == 0 {
		return
	}
	idx, ok := findEdge(successor, target.prefix)
	if !ok {
		return
	}
	n.predecessor.parent.successor = append(successor[:idx], successor[idx+1:]...)
}

// DeletePrefix 删除最长公共前缀的所有有效值 ok 为 false 的时候代表没有删除任何值
func (t *Tree) DeletePrefix(key string) (ok bool) {
	var delCount uint64
	n := t.findLongestPrefixNode(key)
	if n == nil {
		return false
	}

	t.walk(n, func(key string, val any) bool {
		delCount++
		return true
	})

	var parent *node
	if n.predecessor != nil && n.predecessor.parent != nil {
		parent = n.predecessor.parent
	}
	oriPredecessor := n.predecessor
	n.deleteEdge(oriPredecessor)

	if parent != nil {
		parent.mergeChild()
	}
	t.count -= delCount
	return true
}

// mergeChild node 永远是 parent ，如果节点只有一条后继边，并且当前节点没有 leaf 则进行合并
func (n *node) mergeChild() {
	if n.isLeaf() || len(n.successor) > 1 || n.predecessor == nil {
		return
	}
	// 拼接前继边的公共前缀
	n.predecessor.prefix += n.successor[0].prefix
	// 直接继承后继的 leaf
	n.leaf = n.successor[0].getLeaf()
	// 清空后继
	n.successor = nil
}

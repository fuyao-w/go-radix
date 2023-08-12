package go_radix

import (
	. "github.com/fuyao-w/common-util"
	"sort"
	"strings"
)

/*
radix 树实现
*/
type (
	// node 树的节点
	node[T any] struct {
		predecessor *edge[T]   // 节点的前继边
		successor   []*edge[T] // 节点的后继边，按 edge.prefix[0] 从小到大排序
		leaf        *leaf[T]   // 当前节点是否有元素
	}
	// edge 树的边
	edge[T any] struct {
		child  *node[T]
		prefix string // 每个边保存局部的前缀，从顶点到某一叶子节点的局部前缀组成完整的 key
	}
	// Tree 通过树对外提供 API
	Tree[T any] struct {
		root  *node[T]
		count uint64
	}
	leaf[T any] Tuple[string, T]
)

// New 创建一个新的 radix 树
func New[T any]() *Tree[T] {
	return &Tree[T]{root: &node[T]{}}
}

func (e *edge[T]) hasLeaf() bool {
	return e.child.leaf != nil
}

func (e *edge[T]) setLeaf(key string, value T) {
	e.child.setLeaf(key, value)
}
func (e *edge[T]) delLeaf() {
	e.child.leaf = nil
}

func (e *edge[T]) getLeaf() *leaf[T] {
	return e.child.leaf
}
func (n *node[T]) setLeaf(key string, value T) {
	n.leaf = (*leaf[T])(Ptr(BuildTuple(key, value)))
}

func (n *node[T]) isLeaf() bool {
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
func findEdge[T any](s []*edge[T], target string) (idx int, found bool) {
	return sort.Find(len(s), func(i int) int {
		return -strings.Compare(s[i].prefix[:1], target[0:1])
	})
}

// appendEdge 追加一条边，按照 edge.prefix[0] 从小到大的位置进行插入
func (n *node[T]) appendEdge(e *edge[T]) {
	//e.parent = n
	idx, found := findEdge(n.successor, e.prefix)
	if found {
		panic("duplicate edge:" + e.prefix)
	}
	n.successor = append(n.successor[:idx], append([]*edge[T]{e}, n.successor[idx:]...)...)
}

// Insert 插入元素 ok 为 false 代表跟新 ，old 返回之前的 value
func (t *Tree[T]) Insert(key string, value T) (old T, ok bool) {
	curNode := t.root
	search := key

	for {
		// 如果 prefix 是 search 的完全子串，直接终止
		if len(search) == 0 {
			if !curNode.isLeaf() {
				t.count++
				curNode.setLeaf(key, value)
				ok = true
			} else {
				old = curNode.leaf.B
			}
			return
		}
		// 如果 prefix 是 search 的完全子串，直接终止
		idx, found := findEdge(curNode.successor, search)
		if !found {
			next := newEdge[T](search, nil)
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
		search = newPrefix
		if oriPrefix != curEdge.prefix { // prefix 不是 search 子串 (部分重叠)，先拆分，再确定是否增加子节点
			//创建新的边，并且将当前边的 leaf 移动到新边
			next := newEdge(curEdge.prefix[len(oriPrefix):], curEdge.getLeaf())
			next.child.successor = append([]*edge[T](nil), curEdge.child.successor...)
			//next.parent = curEdge.child
			// 将当前边的下游节点后继更新成 next , leaf 置空 ，更新 prefix
			//curEdge.child = next.parent
			curEdge.child.successor = append([]*edge[T](nil), next)
			curEdge.delLeaf()
			curEdge.prefix = oriPrefix
		}
		curNode = curEdge.child
	}
}

// newEdge 创建一条新的边
func newEdge[T any](localPrefix string, leaf *leaf[T]) *edge[T] {
	edge := &edge[T]{prefix: localPrefix}
	edge.child = &node[T]{
		predecessor: edge,
		leaf:        leaf,
	}
	return edge
}

// scan 广度遍历节点 n，walkFn 返回 true 的时候会继续遍历，否则会停止
func (t *Tree[T]) scan(n *node[T], walkFn func(key string, val T) bool) {
	if n == nil {
		return
	}
	var (
		successor, nextSuccessor []*edge[T]
	)
	if n.isLeaf() { // 根节点
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

// Scan 从匹配 prefix 最长前缀的节点开始遍历树，walkFn 返回 true 代表继续遍历，否则终止
func (t *Tree[T]) Scan(prefix string, walkFn func(key string, val T) bool) {
	n, _ := t.findLongestPrefixNode(prefix)
	t.scan(n, walkFn)
}

// findLongestPrefixNode 找到具有最长公共前缀的节点
func (t *Tree[T]) findLongestPrefixNode(prefix string) (curNode, parent *node[T]) {
	curNode = t.root
	search := prefix
	if len(prefix) == 0 {
		return curNode, nil
	}
	var successor = curNode.successor
	for len(successor) > 0 {
		idx, found := findEdge(successor, search)
		if !found {
			break
		}
		curEdge := curNode.successor[idx]
		parent = curNode
		curNode = curEdge.child
		if len(curEdge.prefix) >= len(search) {
			return curNode, parent
		}
		successor = curNode.successor
		search = search[len(curEdge.prefix):]
	}
	return nil, nil
}

// Get 等值查找
func (t *Tree[T]) Get(key string) (val T, exist bool) {
	n, _ := t.findLongestPrefixNode(key)
	if n.isLeaf() {
		return n.leaf.B, true
	}
	return
}

// Len 树的有效元素数量
func (t *Tree[T]) Len() uint64 {
	return t.count
}

// LongestPrefix 返回有效的最长公共前缀
func (t *Tree[T]) LongestPrefix(prefix string) (longestPrefix string, value T, ok bool) {
	n, _ := t.findLongestPrefixNode(prefix)
	if n == nil {
		return
	}
	if n.isLeaf() {
		return n.leaf.A, n.leaf.B, true
	}
	return
}

// Delete 等值删除
func (t *Tree[T]) Delete(key string) (ok bool) {
	n, parent := t.findLongestPrefixNode(key)
	if n == nil {
		return false
	}
	if !n.isLeaf() || n.leaf.A != key {
		return false
	}
	n.leaf = nil
	/*
				删除一个节点的四种情况，
				1. 该节点是中间节点，如果后继只有一条边则将后继合并到当前节点
				例如： a -> b -> c ，删除 b => a -> bc
				2. 该节点是中间节点，如果后继不符合合并条件，但是该节点对于父节点来说符合合并条件，则将该节点合并到父节点中
				例如： a-> b -> c    删除 b => ab -> c
			                ↘  d                  ↘  d
				3. 该节点是中间节点，后继不符合合并条件，该节点对于父节点来说也不符合合并条件，则直接清楚叶子值即可
				例如： a -> b -> c   删除 b => 结构不变
					   ↘    ↘
		                 x   k
			    4. 该节点是叶子节点，直接从父节点中删除和该节点相连的边即可

	*/
	n.mergeChild()
	if parent != nil {
		if len(n.successor) == 0 && !n.isLeaf() {
			parent.deleteEdge(n.predecessor) // 如果是叶子边，可以直接删除
		}

		parent.mergeChild()
	}
	t.count--
	return ok
}

// deleteEdge 从点前节点的后继边删除 target
func (n *node[T]) deleteEdge(target *edge[T]) {
	if target == nil {
		return
	}
	if n.predecessor == nil {
		return
	}
	successor := n.successor
	idx, ok := findEdge(successor, target.prefix)
	if !ok {
		return
	}
	n.successor = append(successor[:idx], successor[idx+1:]...)
}

// DeletePrefix 删除最长公共前缀的所有有效值 ok 为 false 的时候代表没有删除任何值
func (t *Tree[T]) DeletePrefix(key string) (ok bool) {
	var delCount uint64
	n, parent := t.findLongestPrefixNode(key)
	if n == nil {
		return false
	}
	t.scan(n, func(key string, val T) bool {
		delCount++
		return true
	})
	parent.deleteEdge(n.predecessor)
	if parent != nil {
		parent.mergeChild()
	}
	t.count -= delCount
	return true
}

// mergeChild node 永远是 parent ，如果节点只有一条后继边，并且当前节点没有 leaf 则进行合并
func (n *node[T]) mergeChild() {
	// TODO 如果该节点不是叶子节点，并且该节点没有后继，其实可以从前继的边列表中删除该节点所属的边，但是目前没有存储边的前继节点
	if n.isLeaf() || len(n.successor) != 1 || n.predecessor == nil {
		return
	}
	// 拼接前继边的公共前缀
	n.predecessor.prefix += n.successor[0].prefix
	// 直接继承后继的 leaf

	n.leaf = n.successor[0].getLeaf()
	// 清空后继
	n.successor = n.successor[0].child.successor
}

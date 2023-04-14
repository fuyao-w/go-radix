package go_radix

import (
	"testing"
)

func TestInsert(t *testing.T) {
	root := New()
	root.Insert("test", nil)
	root.Insert("tester", nil)
	root.Insert("te", nil)
	root.Insert("tee", nil)
	root.Insert("testerr", nil)
	root.Insert("testerrr", nil)
	t.Log(root)
	root.Scan("", func(key string, val any) bool {
		t.Log(key)
		return true
	})

}

func TestInsert2(t *testing.T) {
	root := New()
	root.Insert(string([]byte("test")), nil)
	root.Insert(string([]byte("toaster")), nil)
	root.Insert(string([]byte("toasting")), nil)
	root.Insert(string([]byte("slow")), 33)
	root.Insert(string([]byte("slowly")), nil)
	root.Insert(string([]byte("slowly")), nil)
	t.Log(root)

	root.Scan("", func(key string, val any) bool {
		t.Log(key)
		return true
	})
	t.Log(root.Len())
	t.Log(root.LongestPrefix("slo"))
}

func TestMatch(t *testing.T) {
	t.Log(longestPrefix("123", "123"))
	t.Log(longestPrefix("123", "12"))
	t.Log(longestPrefix("123", "234"))
	t.Log(longestPrefix("toast", "test"))
	t.Log("123"[:3])
}

func TestFind(t *testing.T) {
	//var slice = []int{1, 2, 3, 4, 5}
	//t.Log(sort.Find(len(slice), func(i int) int {
	//	return 2 - slice[i]
	//}))
	//t.Log(sort.Find(len(slice), func(i int) int {
	//	return 1 - slice[i]
	//}))
	//t.Log(sort.Find(len(slice), func(i int) int {
	//	return -1 - slice[i]
	//}))
	//t.Log(sort.Find(len(slice), func(i int) int {
	//	return 5 - slice[i]
	//}))
	var edges = []*edge{{prefix: "1"}, {prefix: "2"}, {prefix: "3"}, {prefix: "4"}, {prefix: "5"}}
	//var edges = []*edge{{prefix: "1"}, {prefix: "2"}}
	t.Log("1"[0], "2"[0])
	node := &node{successor: edges}
	node.appendEdge(&edge{prefix: "2"})
	node.appendEdge(&edge{prefix: "2"})
	for _, e := range node.successor {
		t.Log(e.prefix)
	}
	t.Log(findEdge(edges, "66"))
}

func TestBlank(t *testing.T) {
	root := New()
	root.Insert("1", 1)
	root.Insert("11", 1)
	root.Insert("111", 1)
	root.Insert("112", 1)
	root.Insert("11234", 1)
	root.Insert("2", 1)
	root.Insert("22", 1)
	root.Insert("", 2)
	root.Scan("112", func(key string, val any) bool {
		t.Log(key, val)
		return true
	})
	t.Log(root.Len())
}

func TestDelete(t *testing.T) {
	root := New()
	root.Insert("test", 1)
	root.Insert("", 1)
	root.Insert("tee", 1)
	root.Delete("tee")
	root.Delete("")
	t.Log(root)
	t.Log(root.Len())
	t.Log(root.Get("test"))
}

func TestDeletePrefix(t *testing.T) {
	root := New()
	root.Insert("test", 1)
	root.Insert("", 1)
	root.Insert("tee", 1)
	root.DeletePrefix("te")
	t.Log(root)
	t.Log(root.Len())
}

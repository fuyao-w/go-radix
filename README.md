# go-radix

## radix 树 go 实现，支持泛型

### 单测覆盖率 89.6%


- 用例

```go

func main() {
	tree := New[int]()
	tree.Insert("abc", 1)
	tree.Insert("abac", 2)
	tree.Insert("d", 3)
	tree.Delete("d")
	val, ok := tree.Get("abc")
	if !ok {
		log.Println("key abc not exist")
	} else {
		log.Println("value ", val)
	}

	prefix, value, ok := tree.LongestPrefix("a")
	if !ok {
		log.Println("key a not have common prefix")
	} else {
		log.Println("prefix ", prefix, "value", value)
	}
	// 遍历所有节点
	tree.Scan("", func(key string, val int) bool {
		log.Println("scan", "key", key, "val", val)
		return true
	})
}

```
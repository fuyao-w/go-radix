package go_radix

import (
	"github.com/gookit/goutil/dump"
	"testing"
)

func TestInsert(t *testing.T) {
	tests := []struct {
		key string
		val string
	}{
		{
			"test", "test",
		}, {
			"t", "t",
		}, {
			"toasting", "toasting",
		}, {
			"xy", "xy",
		}, {
			"xy", "xy",
		}, {
			"tester", "tester",
		}, {
			"te", "te",
		}, {
			"tee", "tee",
		}, {
			"testerr", "testerr",
		}, {
			"testerr", "testerr",
		}, {
			"1", "1",
		}, {
			"2", "2",
		}, {
			"11", "11",
		}, {
			"22", "22",
		}, {
			"111", "111",
		}, {
			"222", "222",
		},
	}
	root := New[string]()
	t.Run("", func(t *testing.T) {
		for _, test := range tests {
			root.Insert(test.key, test.val)
		}
		for _, test := range tests {
			val, ok := root.Get(test.key)
			if !ok {
				t.Fatalf("get key not found :%s", test.key)
			}
			if val != test.val {
				t.Fatalf("get key not found :%s ,res :%s", test.key, val)
			}
		}
	})
	t.Run("length", func(t *testing.T) {
		if root.Len() != uint64(len(tests)-2) {
			t.Fatalf("lenth not equal :%d ,res :%d", root.Len(), len(tests)-2)
		}
	})

}

func Test_Delete(t *testing.T) {
	tests := []struct {
		key string
		val string
	}{
		{
			"test", "test",
		}, {
			"t", "t",
		}, {
			"toasting", "toasting",
		},
		{
			"xy", "xy",
		},
		{
			"tester", "tester",
		}, {
			"te", "te",
		}, {
			"tee", "tee",
		}, {
			"testerr", "testerr",
		}, {
			"testerrr", "testerrr",
		}, {
			"testerar", "testerar",
		},
		{
			"1", "1",
		},
		{
			"2", "2",
		}, {
			"11", "11",
		}, {
			"22", "22",
		}, {
			"111", "111",
		}, {
			"222", "222",
		},
	}
	root := New[string]()
	for _, test := range tests {
		root.Insert(test.key, test.val)
	}
	t.Run("length", func(t *testing.T) {
		if root.Len() != uint64(len(tests)) {
			t.Fatal("lenth not equal ")
		}
	})
	for _, test := range tests {
		root.Delete(test.key)
	}
	t.Run("length", func(t *testing.T) {
		if root.Len() != 0 {
			t.Fatalf("lenth not zero :%d ", root.Len())
		}
	})
}

func TestInsert2(t *testing.T) {
	root := New[any]()
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

	n := &node[any]{successor: nil}
	for i := 97 + 25; i >= 97; i-- {
		n.appendEdge(newEdge[any](string([]byte{byte(i)}), nil))
	}

	//var edges = []*edge{{prefix: "1"}, {prefix: "2"}, {prefix: "3"}, {prefix: "4"}, {prefix: "5"}}
	//var edges = []*edge{{prefix: "1"}, {prefix: "2"}}
	//t.Log("1"[0], "2"[0])
	//
	//node.appendEdge(&edge{prefix: "2"})
	//node.appendEdge(&edge{prefix: "2"})
	for _, e := range n.successor {
		t.Log(e.prefix)
	}
	n = &node[any]{}
	n.appendEdge(newEdge[any]("4", nil))
	n.appendEdge(newEdge[any]("1", nil))
	n.appendEdge(newEdge[any]("3", nil))
	n.appendEdge(newEdge[any]("5", nil))
	n.appendEdge(newEdge[any]("2", nil))
	for _, e := range n.successor {
		t.Log(e.prefix)
	}
	//t.Log(findEdge(edges, "66"))
}

func TestBlank(t *testing.T) {
	root := New[int]()
	root.Insert("1", 1)
	root.Insert("11", 1)
	root.Insert("111", 1)
	root.Insert("112", 1)
	root.Insert("113", 1)
	root.Insert("11234", 1)
	root.Insert("2", 1)
	root.Insert("22", 1)
	root.Insert("", 2)
	root.Scan("112", func(key string, val int) bool {
		t.Log(key, val)
		return true
	})
	t.Log(root.Len())
}

func TestDelete(t *testing.T) {
	root := New[int]()
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
	root := New[int]()
	root.Insert("1", 1)
	root.Insert("11", 1)
	root.Insert("111", 1)
	root.Insert("112", 1)
	root.Insert("113", 1)
	root.Insert("11234", 1)
	root.Insert("2", 1)
	root.Insert("22", 1)
	root.DeletePrefix("112")
	t.Log(root)
	if root.Len() != 6 {
		t.Fatal("DeletePrefix fail")
	}
}

func Test_longestPrefix(t *testing.T) {
	tests := []struct {
		item, target string
		wantIdx      int
	}{
		{
			"abc",
			"ab",
			2,
		},
		{
			"cdd",
			"ab",
			0,
		},
		{
			"123",
			"123",
			3,
		},
		{
			"1",
			"123",
			1,
		},
		{
			"",
			"123",
			0,
		},
		{
			"",
			"",
			0,
		},
		{
			"1",
			"",
			0,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if gotIdx := longestPrefix(tt.item, tt.target); gotIdx != tt.wantIdx {
				t.Errorf("longestPrefix() = %v, want %v", gotIdx, tt.wantIdx)
			}
		})
	}
}

func Test_findEdge(t *testing.T) {
	newEdge := func(prefix string) *edge[any] {
		return newEdge[any](prefix, nil)
	}
	type args struct {
		s      []*edge[any]
		target string
	}
	tests := []struct {
		//name      string
		args      args
		wantIdx   int
		wantFound bool
	}{
		{
			args: args{
				s:      []*edge[any]{newEdge("11"), newEdge("22"), newEdge("33")},
				target: "23445",
			},
			wantIdx:   1,
			wantFound: true,
		}, {
			args: args{
				s:      []*edge[any]{newEdge(""), newEdge("22"), newEdge("33")},
				target: "5",
			},
			wantIdx:   3,
			wantFound: false,
		}, {
			args: args{
				s:      []*edge[any]{newEdge("1"), newEdge("22"), newEdge("33")},
				target: "-",
			},
			wantIdx:   0,
			wantFound: false,
		}, {
			args: args{
				s:      []*edge[any]{newEdge("1"), newEdge("22"), newEdge("33")},
				target: "223",
			},
			wantIdx:   1,
			wantFound: true,
		}, {
			args: args{
				s:      []*edge[any]{newEdge("1"), newEdge("22"), newEdge("33")},
				target: "3",
			},
			wantIdx:   2,
			wantFound: true,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gotIdx, gotFound := findEdge(tt.args.s, tt.args.target)
			if gotIdx != tt.wantIdx {
				t.Errorf("findEdge() gotIdx = %v, want %v", gotIdx, tt.wantIdx)
			}
			if gotFound != tt.wantFound {
				t.Errorf("findEdge() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func Test_appendEdge(t *testing.T) {
	e := func(prefix string) *edge[any] {
		return newEdge[any](prefix, nil)
	}

	tests := []struct {
		slice, target []string
	}{
		{
			slice:  []string{"1", "44", "2", "543", "3"},
			target: []string{"1", "2", "3", "44", "543"},
		},
		{
			slice:  []string{"e23432", "dsf", "cc", "bs45", "as"},
			target: []string{"as", "bs45", "cc", "dsf", "e23432"},
		}, {
			slice:  []string{"10", "9", "8", "7", "6"},
			target: []string{"10", "6", "7", "8", "9"},
		}, {
			slice:  []string{"9", "110", "8", "7", "6"},
			target: []string{"110", "6", "7", "8", "9"},
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			n := &node[any]{}
			for _, s := range test.slice {
				n.appendEdge(e(s))
			}
			if len(n.successor) != len(test.target) {
				dump.Println(n.successor)
				t.Errorf("appendEdge() length not equal ,want :%d ,result :%d", len(test.target), len(n.successor))
			}
			for i, e := range n.successor {
				if e.prefix != test.target[i] {
					t.Errorf(dump.Format("appendEdge() not equal :%", n.successor))
				}
			}
		})
	}

}

package evaluator

import "testing"

func TestScope(t *testing.T) {
	tables := []struct {
		fn       func() OBJECT
		expected OBJECT
	}{
		{
			func() OBJECT {
				root := CreateScope()
				root.Init("a", Int(1))
				return root.Get("a")
			},
			Int(1),
		},
		{
			func() OBJECT {
				root := CreateScope()
				root.Init("a", String("Hello"))
				root = root.Child()
				root.SetOrInit("a", String("World"))
				root = root.prev
				return root.Get("a")
			},
			String("World"),
		},
		{
			func() OBJECT {
				root := CreateScope()
				root.SetOrInit("a", String("Hello"))
				root = root.Child()
				root.Init("a", String("World"))
				root = root.prev
				return root.Get("a")
			},
			String("Hello"),
		},
		{
			func() OBJECT {
				root := CreateScope()
				root.Init("b", Float(1.1))
				root.Init("b", Float(1.2))
				root.Init("b", Float(1.3))
				return root.Get("b")
			},
			Float(1.3),
		},
		{
			func() OBJECT {
				return CreateScope().Get("undefined")
			},
			Null_,
		},
		{
			func() OBJECT {
				root := CreateScope()
				root.Init("t", Int(1))
				root = root.Child().Child()
				root.SetOrInit("t", Int(2))
				root.Init("t", Null_)
				root.SetOrInit("t", Int(4))
				root = root.prev.prev
				return root.Get("t")
			},
			Int(2),
		},
	}

	for _, table := range tables {
		result := table.fn()
		if table.expected != result {
			t.Errorf("Expected: %s, found: %s",
				table.expected.ToString(), result.ToString())
		}
	}
}

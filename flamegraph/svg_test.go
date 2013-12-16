package flamegraph

import (
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	calls := [][]Call{
		{{"a", "a"}, {"a", "b"}},
		{{"b", "a"}, {"a", "a"}},
		{{"a", "a"}},
		{{"a", "b"}, {"a", "a"}},
	}

	expected := [][]Call{
		{{"a", "a"}},
		{{"a", "a"}, {"a", "b"}},
		{{"a", "b"}, {"a", "a"}},
		{{"b", "a"}, {"a", "a"}},
	}

	s := make(stacks, len(calls))
	for i := range calls {
		s[i] = &Stack{Calls: calls[i]}
	}

	sort.Sort(s)

	for i, stack := range s {
		if len(stack.Calls) != len(expected[i]) {
			t.Errorf("Len %d != %d; got %v, want %v", len(stack.Calls), len(expected[i]),
				stack.Calls, expected)
		}
	}
}

func TestDo(t *testing.T) {
	calls := [][]Call{
		{{"a", "a"}, {"a", "b"}},
		{{"b", "a"}, {"a", "a"}},
		{{"a", "a"}},
		{{"a", "b"}, {"a", "a"}},
	}

	tr := make(traces, len(calls))
	for i, c := range calls {
		tr[i] = trace{&Stack{Calls: c}, i}
	}

	nodes := makeNodes(tr)
	//TODO test that nodes is correct.
	t.Log(nodes)
}

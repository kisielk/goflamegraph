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

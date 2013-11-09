package flamegraph

import (
	"sort"
	"strconv"
	"strings"
)

func foldStack(s Stack) string {
	calls := make([]string, len(s.Calls))
	for i, c := range s.Calls {
		calls[i] = c.Source + "`" + c.Func
	}
	return strings.Join(calls, ";")
}

func foldStacks(stacks []*Stack) []string {
	counts := make(map[string]int)
	for _, s := range stacks {
		folded := foldStack(*s)
		counts[folded]++
	}

	lines := make([]string, 0, len(counts))
	for line, count := range counts {
		lines = append(lines, line+" "+strconv.Itoa(count))
	}

	sort.Strings(lines)
	return lines
}

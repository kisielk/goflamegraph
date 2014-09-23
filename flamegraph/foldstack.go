package flamegraph

import (
	"sort"
	"strconv"
	"strings"
)

func foldStack(s Stack, includeSource bool) string {
	calls := make([]string, len(s.Calls))
	for i, c := range s.Calls {
		call := c.Func
		if includeSource {
			call = c.Source + "`" + c.Func
		}
		calls[i] = call
	}
	return strings.Join(calls, ";")
}

func FoldStacks(stacks []*Stack, includeSource bool) []string {
	counts := make(map[string]int)
	for _, s := range stacks {
		folded := foldStack(*s, includeSource)
		counts[folded]++
	}

	lines := make([]string, 0, len(counts))
	for line, count := range counts {
		lines = append(lines, line+" "+strconv.Itoa(count))
	}

	sort.Strings(lines)
	return lines
}

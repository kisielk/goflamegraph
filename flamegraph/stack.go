// The code in this file is from https://code.google.com/p/rog-go/source/browse/cmd/stackgraph/stackgraph.go
package flamegraph

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

type Call struct {
	Func   string
	Source string
}

type Stack struct {
	Goroutine int
	State     string
	Calls     []Call
}

func ParseStacks(r io.Reader) ([]*Stack, error) {
	var stacks []*Stack
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		stack := &Stack{}
		if n, _ := fmt.Sscanf(line, "goroutine %d [%s]", &stack.Goroutine, &stack.State); n != 2 {
			continue
		}
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				// empty line signifies end of a stack
				break
			}
			if strings.Contains(line, "  ") {
				// Looks like a register dump.
				// TODO better heuristic here.
				continue
			}
			if strings.HasSuffix(line, ")") {
				if i := strings.LastIndex(line, "("); i > 0 {
					line = line[0:i]
				}
			}
			line = strings.TrimPrefix(line, "created by ")
			call := Call{Func: line}
			if !scanner.Scan() {
				break
			}
			line = scanner.Text()
			if strings.HasPrefix(line, "\t") {
				line = strings.TrimPrefix(line, "\t")
				if i := strings.LastIndex(line, " +"); i >= 0 {
					line = line[0:i]
				}
				call.Source = line
			}
			stack.Calls = append(stack.Calls, call)
		}
		if len(stack.Calls) > 0 {
			stacks = append(stacks, stack)
		}
	}
	return stacks, nil
}

func foldStack(s Stack) string {
	calls := make([]string, len(s.Calls))
	for i, c := range s.Calls {
		calls[i] = c.Source + "`" + c.Func
	}
	return strings.Join(calls, ";")
}

func FoldStacks(stacks []*Stack) []string {
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

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

// stackLess returns true if stack a is "less than" stack b.
func stackLess(a, b *Stack) bool {
	ca := a.Calls
	cb := b.Calls
	for x := 0; x < len(ca) && x < len(cb); x++ {
		cas := ca[x].Source
		cbs := cb[x].Source
		if cas < cbs {
			return true
		} else if cas > cbs {
			return false
		}

		caf := ca[x].Func
		cbf := cb[x].Func
		if caf < cbf {
			return true
		} else if caf > cbf {
			return false
		}

		// The functions at the current level of the stacks are equal, continue.
	}

	// If we reach here, then all the functions have been equal at all levels inspected.
	// Check if one of the stacks is smaller than the other.
	return len(ca) < len(cb)
}

type stacks []*Stack

func (s stacks) Len() int {
	return len(s)
}

func (s stacks) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s stacks) Less(i, j int) bool {
	return stackLess(s[i], s[j])
}

type trace struct {
	stack   *Stack
	samples int
}

type traces []trace

func (t traces) Len() int {
	return len(t)
}

func (t traces) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t traces) Less(i, j int) bool {
	return stackLess(t[i].stack, t[j].stack)
}

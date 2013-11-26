package flamegraph

import (
	"fmt"
	"sort"
)

const svgTemplate = `<?xml version="1.0" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg version="1.1" width="{{width}}" height="{{height}}" onload="init(evt)" viewBox="0 0 {{width}} {{height}}" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
<defs>
        <linearGradient id="background" y1="0" y2="1" x1="0" x2="0" >
                <stop stop-color="{{bgColor1}}" offset="5%" />
                <stop stop-color="{{bgColor2}}" offset="95%" />
        </linearGradient>
</defs>
<style type="text/css">
        .func_g:hover { stroke:black; stroke-width:0.5; }
</style>
<script type="text/ecmascript">
<![CDATA[
        var details;
        function init(evt) { details = document.getElementById("details").firstChild; }
        function s(info) { details.nodeValue = "$nametype " + info; }
        function c() { details.nodeValue = ' '; }
]]>
</script>
<rect x="0" y="0" width="{{width}}" height="{{height}}" fill="{{background}}" />

</svg>
`

type color struct{ r, g, b int }

func (c color) String() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", c.r, c.g, c.b)
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

func do(t traces) {
	sort.Sort(t)
	var prev *Stack
	var totalSamples int
	for _, trace := range t {
		if trace.samples <= 0 {
			continue
		}

		flow(prev, trace.stack, totalSamples)
		prev = trace.stack
		totalSamples += trace.samples
	}
}

func flow(prev, curr *Stack, totalSamples int) {

}

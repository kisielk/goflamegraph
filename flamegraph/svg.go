package flamegraph

import (
	"html/template"
	"io"
	"sort"
)

const svg = `
<?xml version="1.0" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg version="1.1" width="{{.Width}}" height="{{.Height}}" onload="init(evt)" viewBox="0 0 {{.Width}} {{.Height}}" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
<defs >
	<linearGradient id="background" y1="0" y2="1" x1="0" x2="0" >
		<stop stop-color="{{.BgColor1}}" offset="5%" />
		<stop stop-color="{{.BgColor2}}" offset="95%" />
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
<rect x="0" y="0" width="{{.Width}}" height="{{.Height}}" fill="url(#background)"/>
<text text-anchor="middle" x="{{.Width | div 2}}" y="{{.FontSize | mulFloat 2}}" font-size="{{.FontSize}}" font-family="{{.FontFamily}}" fill="rgb(0,0,0)">
{{.Title}}
</text>
{{range .Funcs}}
<g class="func_g" onmouseover="s('{{.Info}}')" onmouseout="c()">
<title>{{.Info}}</title><rect x="144.1" y="209" width="1019.1" height="15.0" fill="rgb(215,143,50)" rx="2" ry="2" />
<text text-anchor="" x="147.090909090909" y="219.5" font-size="12" font-family="Verdana" fill="rgb(0,0,0)">{{.Name}}</text>
</g>
{{end}}
</svg>
`

var svgTemplate = template.Must(template.New("SVG").Funcs(funcMap).Parse(svg))

type Args struct {
	Width, Height      int
	BgColor1, BgColor2 string
	FontSize           float64
	FontFamily         string
	Title              string
	Funcs              []Func
}

type Func struct {
	Info, Name   string
	Etime, Stime float64
}

var funcMap = template.FuncMap{
	"mulFloat": func(a, b float64) float64 { return a * b },
	"mul":      func(a, b int) int { return a * b },
	"div": func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a / b
	},
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

func RenderSVG(w io.Writer) error {
	args := Args{
		Width:      970,
		Height:     640,
		BgColor1:   "black",
		BgColor2:   "blue",
		FontFamily: "Helvetica",
		FontSize:   10,
		Title:      "Something",
	}
	return svgTemplate.Execute(w, args)
}

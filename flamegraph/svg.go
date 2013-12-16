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
{{range $func, $depth := .Nodes}}
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
	Nodes              map[node]int
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

const minWidthTime = 1

func makeNodes(t traces) map[node]int {
	tmp := make(map[node]int)
	nodes := make(map[node]int)

	sort.Sort(t)
	var prev []Call
	var totalSamples int
	for _, trace := range t {
		if trace.samples <= 0 {
			continue
		}

		curr := trace.stack.Calls
		flow(tmp, nodes, prev, curr, totalSamples)
		prev = curr
		totalSamples += trace.samples
	}

	// Prune nodes that are too narrow, find maximum depth.
	var maxDepth int
	for node, startTime := range nodes {
		if node.EndTime-startTime < minWidthTime {
			delete(nodes, node)
			continue
		}
		if node.Depth > maxDepth {
			maxDepth = node.Depth
		}
	}

	return nodes
}

type node struct {
	Call    Call
	Depth   int
	EndTime int
}

func flow(tmp map[node]int, nodes map[node]int, prev, this []Call, totalSamples int) {
	var same int
	var i int
	for i = 0; i < len(prev) && i < len(this); i++ {
		if prev[i] != this[i] {
			break
		}
	}
	same = i

	for i := len(prev) - 1; i >= same; i-- {
		k := node{Call: prev[i], Depth: i}
		nodes[node{k.Call, k.Depth, totalSamples}] = tmp[k]
		delete(tmp, k)
	}

	for i := same; i < len(this); i++ {
		k := node{Call: this[i], Depth: i}
		tmp[k] = totalSamples
	}
}

func RenderSVG(w io.Writer, t traces) error {
	nodes := makeNodes(t)
	args := Args{
		Width:      970,
		Height:     640,
		BgColor1:   "black",
		BgColor2:   "blue",
		FontFamily: "Helvetica",
		FontSize:   10,
		Title:      "Something",
		Nodes:      nodes,
	}
	return svgTemplate.Execute(w, args)
}

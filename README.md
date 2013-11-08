# goflamegraph

goflamegraph helps generate flamegraphs from Go stack traces.

## Usage

First install this tool:

    go get github.com/kisielk/goflamegraph

You will also need flamegraph.pl from https://github.com/brendangregg/FlameGraph

Then generate a bunch of stack traces. There are several facilities to do this in Go.

You can periodically run something like:

    prof := pprof.Lookup("goroutine")
    prof.WriteTo(os.Stderr, 2)

inside your program.

Another option is to install the HTTP handler from net/http/pprof and periodically
sample the /debug/pprof/goroutine?debug=2 URL.

Note that in both cases you must ensure that each dump is separated by an extra
newline or else the output will apear to be incorrect.

Once you have your stack traces, you can filter them through goflamegraph and flamgraph.pl:

    goflamegraph < stacks.txt | flamegraph.pl > out.svg

Then view the SVG file in a web browser to get the interactive output.

## TODO

Some future developments:

  * Rewrite SVG rendering in Go.
  * Automate stack collection.
  * Provide an HTTP handler to view a flamegraph of a running program.

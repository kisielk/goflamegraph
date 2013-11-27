package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kisielk/goflamegraph/flamegraph"
)

var html = flag.Bool("html", false, "Generate HTML instead of output for flamegraph.pl")

func main() {
	flag.Parse()

	stacks, err := flamegraph.ParseStacks(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	if *html {
		err := flamegraph.RenderSVG(os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		for _, line := range flamegraph.FoldStacks(stacks) {
			fmt.Println(line)
		}
	}
}

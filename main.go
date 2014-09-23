package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kisielk/goflamegraph/flamegraph"
)

var includeSource = flag.Bool("s", false, "Show source paths")

func main() {
	stacks, err := flamegraph.ParseStacks(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range flamegraph.FoldStacks(stacks, *includeSource) {
		fmt.Println(line)
	}
}

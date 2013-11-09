package main

import (
	"fmt"
	"github.com/kisielk/goflamegraph/flamegraph"
	"log"
	"os"
)

func main() {
	stacks, err := flamegraph.ParseStacks(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range flamegraph.FoldStacks(stacks) {
		fmt.Println(line)
	}
}

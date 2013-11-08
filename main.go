package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	stacks, err := parseStacks(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range foldStacks(stacks) {
		fmt.Println(line)
	}
}

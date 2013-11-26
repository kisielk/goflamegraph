package flamegraph

import (
	"bytes"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

func init() {
	http.Handle("/debug/flamegraph", http.HandlerFunc(FlameGraph))
}

func FlameGraph(w http.ResponseWriter, r *http.Request) {
	sec, _ := strconv.ParseInt(r.FormValue("seconds"), 10, 64)
	if sec == 0 {
		sec = 30
	}

	freq, _ := strconv.ParseInt(r.FormValue("frequency"), 10, 64)
	if freq == 0 {
		freq = 997
	}

	size := 4096
	ticker := time.NewTicker(1 * time.Second / time.Duration(freq))
	timeout := time.NewTimer(time.Duration(sec) * time.Second)
	var stacks [][]byte
loop:
	for {
		select {
		case <-ticker.C:
		tick:
			buf := make([]byte, size)
			n := runtime.Stack(buf, true)
			if n == size {
				size *= 2
				goto tick
			}
			stacks = append(stacks, buf[:n])
		case <-timeout.C:
			break loop
		}
	}

	var allStacks []*Stack
	for _, s := range stacks {
		currentStacks, err := ParseStacks(bytes.NewReader(s))
		if err != nil {
			log.Println("goflamegraph:", err)
			continue
		}
		allStacks = append(allStacks, currentStacks...)
	}

	// TODO: Render SVG
}

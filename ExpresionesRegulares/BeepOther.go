package main

import (
	"flag"
	"fmt"
)

var (
	times int
)

func init() {
	flag.IntVar(&times, "times", 3, "Number of times to beep.")
	flag.Parse()
}

func main() {
	for i := 0; i < times; i++ {
		fmt.Print("\a")
	}
}

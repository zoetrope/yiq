package main

import (
	"fmt"
	"os"

	"github.com/fiatjaf/jiq"
)

func main() {
	content := os.Stdin

	outputquery := false
	jqargs := os.Args[1:]
	for i, arg := range os.Args[1:] {
		if arg == "-q" {
			outputquery = true
			jqargs = os.Args[1 : i-1]
			jqargs = append(jqargs, os.Args[i:]...)
			break
		}
	}

	e := jiq.NewEngine(content, jqargs)
	os.Exit(run(e, outputquery))
}

func run(e *jiq.Engine, outputquery bool) int {
	result := e.Run()
	if result.Err != nil {
		return 2
	}
	if outputquery {
		fmt.Printf("%s", result.Qs)
	} else {
		fmt.Printf("%s", result.Content)
	}
	return 0
}

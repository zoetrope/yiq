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

func run(e jiq.EngineInterface, outputquery bool) int {
	result := e.Run()
	if result.GetError() != nil {
		return 2
	}
	if outputquery {
		fmt.Printf("%s", result.GetQueryString())
	} else {
		fmt.Printf("%s", result.GetContent())
	}
	return 0
}

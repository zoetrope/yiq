package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fiatjaf/jiq"
)

func main() {
	content := os.Stdin

	var query bool

	flag.BoolVar(&query, "q", false, "output query")
	flag.Parse()

	e := jiq.NewEngine(content)
	os.Exit(run(e, query))
}

func run(e jiq.EngineInterface, query bool) int {

	result := e.Run()
	if result.GetError() != nil {
		return 2
	}
	if query {
		fmt.Printf("%s", result.GetQueryString())
	} else {
		fmt.Printf("%s", result.GetContent())
	}
	return 0
}

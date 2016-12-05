package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fiatjaf/jiq"
)

func main() {
	content := os.Stdin

	initialquery := "."
	outputquery := false
	jqargs := os.Args[1:]
	for i, arg := range os.Args[1:] {
		i = i + 1
		if arg == "-q" {
			outputquery = true
			jqargs = os.Args[1:i]
			jqargs = append(jqargs, os.Args[i+1:]...)
			break
		} else if arg[0] != '-' {
			log.Print(i, " ", arg)
			initialquery = arg
			jqargs = os.Args[1:i]
			jqargs = append(jqargs, os.Args[i+1:]...)
			break
		}
	}

	e := jiq.NewEngine(content, jqargs, initialquery)
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

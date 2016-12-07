package main

import (
	"fmt"
	"os"

	"github.com/fiatjaf/jiq"
)

func main() {
	content := os.Stdin

	if len(os.Args) > 1 && os.Args[1] == "--help" {
		fmt.Print(`jiq - interactive commandline JSON processor
Usage: <json string> | jiq [options] [initial filter]

    jiq is a tool that allows you to play with jq filters interactively
    acting direcly on a JSON source of your choice, given through STDIN.
    For all the details about which filters you can use to transform your
    JSON string, see jq(1) manpage or https://stedolan.github.io/jq

    jiq supports all command line arguments jq supports, plus
     -q         will print the ending filter to STDOUT, instead of
                printing the resulting filtered JSON, the default.
     --help     prints this help message.
      `)
		os.Exit(0)
	}

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

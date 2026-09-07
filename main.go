package main

import (
	"budgetpipe/cmd"
	_ "embed"
)

//go:embed assets/budget-template.xlsx
var templateBytes []byte

func main() {
	if err := cmd.Execute(templateBytes); err != nil {
		panic(err)
	}
}

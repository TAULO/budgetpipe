package main

import (
	_ "embed"

	"github.com/TAULO/budgetpipe/cmd"
)

//go:embed assets/budget-template.xlsx
var templateBytes []byte

func main() {
	if err := cmd.Execute(templateBytes); err != nil {
		panic(err)
	}
}

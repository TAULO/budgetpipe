package main

import "budgetpipe/cmd"

func main() {
	err := cmd.Execute()

	if err != nil {
		panic(err)
	}
}

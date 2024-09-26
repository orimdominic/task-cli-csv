package main

import (
	pkg "github.com/orimdominic/cli-tasks-csv/internal"
)

func main() {
	pkg.Execute("./tasks.csv")
}

// TODO
// Ask user via Readline if they want to update already updated task completion date

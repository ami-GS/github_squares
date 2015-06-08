package main

import (
	. "github.com/ami-GS/github_squares"
	"os"
)

func main() {
	if len(os.Args) == 2 {
		userName := os.Args[1]
		ShowSquare(userName)
	}
}

# github_squares
Code for getting GitHub commit streak squares

## Example
```
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
```

## Usage
```
go run ./example/show_github_squares.go USERNAME
```

### License
The MIT License (MIT) Copyright (c) 2015 ami-GS
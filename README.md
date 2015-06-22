# github_squares
Code for getting GitHub commit streak squares

## Example
```go
package main

import (
	. "github.com/ami-GS/github_squares"
)

func main() {
	ShowSquare()
}
```

## Usage
```
go run ./example/show_github_squares.go USERNAME
```

## Result
![alt tag](https://raw.github.com/ami-GS/github_squares/master/image/example.png)

## TODO
* Retry when scraping failed

### License
The MIT License (MIT) Copyright (c) 2015 ami-GS
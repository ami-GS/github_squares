package github_squares

import (
	"github.com/ami-GS/soac/go"
)

var colorMap map[string]byte = map[string]byte{
	"#d6e685": 156,
	"#8cc665": 112,
	"#44a340": 34,
	"#1e6823": 22,
	"#eeeeee": 237,
}

var Changer *soac.Changer

func init() {
	Changer = soac.NewChanger()
}

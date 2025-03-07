package main

import (
	"github.com/adamdecaf/duperrormsg/duperrormsg"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(duperrormsg.Analyzer)
}

package duperrormsg_test

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/adamdecaf/duperrormsg/duperrormsg"
)

func TestAll(t *testing.T) {
	wd, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, wd, duperrormsg.Analyzer, "tests")
}

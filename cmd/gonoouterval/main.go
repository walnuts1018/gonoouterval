package main

import (
	"github.com/motemen/go-statictools/noouterval"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(noouterval.Analyzer)
}

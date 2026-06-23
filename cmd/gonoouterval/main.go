package main

import (
	"github.com/walnuts1018/noouterval/noouterval"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(noouterval.Analyzer)
}

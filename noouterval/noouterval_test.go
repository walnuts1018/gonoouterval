package noouterval

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testData := analysistest.TestData()

	typePath = "test1.Conn"
	analysistest.Run(t, testData, Analyzer, "test1")
}

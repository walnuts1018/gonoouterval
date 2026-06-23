package noouterval

import (
	"testing"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testData := analysistest.TestData()

	analysistest.Run(t, testData, NewAnalyzer(Settings{
		Type: "testdata.Execer",
	}), ".")
}

func TestPlugin(t *testing.T) {
	newPlugin, err := register.GetPlugin("gonoouterval")
	if err != nil {
		t.Fatal(err)
	}

	plugin, err := newPlugin(map[string]any{
		"type": "testdata.Execer",
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := plugin.GetLoadMode(); got != register.LoadModeTypesInfo {
		t.Fatalf("load mode = %q, want %q", got, register.LoadModeTypesInfo)
	}

	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		t.Fatal(err)
	}
	if len(analyzers) != 1 {
		t.Fatalf("len(analyzers) = %d, want 1", len(analyzers))
	}
	if got := analyzers[0].Name; got != "gonoouterval" {
		t.Fatalf("analyzer name = %q, want %q", got, "gonoouterval")
	}
}

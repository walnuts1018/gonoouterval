// Package noouterval finds uses of values from the outer scope,
// despite values of same type is defined in the inner scope.
package noouterval

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "check for uses of values from the outer scope despite values of same type in the inner scope"

// Settings is the golangci-lint module plugin configuration.
type Settings struct {
	Type string `json:"type"`
}

// Analyzer is the default analyzer used by singlechecker.
var Analyzer = NewAnalyzer(Settings{})

func init() {
	register.Plugin("gonoouterval", newPlugin)
}

// NewAnalyzer returns a gonoouterval analyzer configured for the target type.
func NewAnalyzer(settings Settings) *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name: "gonoouterval",
		Doc:  doc,
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return run(pass, settings.Type)
		},
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}

	a.Flags.StringVar(&settings.Type, "type", settings.Type, "`path/to/pkg.type` to check for")

	return a
}

type plugin struct {
	settings Settings
}

func newPlugin(rawSettings any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[Settings](rawSettings)
	if err != nil {
		return nil, err
	}
	if settings.Type == "" {
		return nil, fmt.Errorf("type is required")
	}

	return &plugin{settings: settings}, nil
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{NewAnalyzer(p.settings)}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

func lookupType(pkg *types.Package, typePath string) (types.Object, error) {
	pos := strings.LastIndex(typePath, ".")
	if pos == -1 {
		return nil, fmt.Errorf("invalid type path: %s", typePath)
	}

	pkgPath := typePath[:pos]
	typeName := typePath[pos+1:]

	visited := map[*types.Package]bool{}
	queue := []*types.Package{pkg}
	for len(queue) > 0 {
		pkg := queue[0]
		queue = queue[1:]

		if visited[pkg] {
			continue
		}
		visited[pkg] = true

		queue = append(queue, pkg.Imports()...)

		if pkg.Path() != pkgPath {
			continue
		}
		typ := pkg.Scope().Lookup(typeName)
		if typ != nil {
			return typ, nil
		}
	}

	return nil, nil
}

func assinableTo(t1, t2 types.Type) bool {
	if bt, ok := t1.(*types.Basic); ok && bt.Kind() == types.Invalid {
		return false
	}
	return types.AssignableTo(t1, t2)
}

func baseIdent(node ast.Node) *ast.Ident {
	switch node := node.(type) {
	case *ast.Ident:
		return node
	case *ast.SelectorExpr:
		return baseIdent(node.X)
	default:
		return nil
	}
}

func run(pass *analysis.Pass, typePath string) (interface{}, error) {
	targetType, err := lookupType(pass.Pkg, typePath)
	if err != nil {
		return nil, err
	}
	if targetType == nil {
		// fail silently
		return nil, nil
	}

	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	ins.WithStack([]ast.Node{(*ast.Ident)(nil)}, func(node ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return false
		}

		foundIdent := node.(*ast.Ident)
		foundIdentObj, ok := pass.TypesInfo.Uses[foundIdent].(*types.Var)
		if !ok {
			return true
		}

		if assinableTo(foundIdentObj.Type(), targetType.Type()) {
			// node is a variable of the targetType
		} else {
			return true
		}

		var (
			foundExpr  ast.Expr = foundIdent
			foundScope          = foundIdentObj.Parent()
		)

		if foundIdentObj.IsField() {
			for i := range stack {
				if expr, ok := stack[len(stack)-1-i].(*ast.SelectorExpr); ok && expr.Sel == foundIdent {
					if base := baseIdent(expr); base != nil {
						foundExpr = expr
						foundScope = pass.TypesInfo.Uses[base].Parent()
						break
					}
					return true
				}
				if kv, ok := stack[len(stack)-1-i].(*ast.KeyValueExpr); ok && kv.Key == foundIdent {
					// node is K in Struct{ K: ... }
					return true
				}
			}
		}

		thisScope := pass.Pkg.Scope().Innermost(foundIdent.Pos())
		if foundScope == thisScope {
			return true
		}

		for thisScope != nil && thisScope != foundScope && thisScope != types.Universe {
			for _, name := range thisScope.Names() {
				objInScope := thisScope.Lookup(name)
				if !assinableTo(objInScope.Type(), targetType.Type()) {
					continue
				}

				pass.Report(
					analysis.Diagnostic{
						Pos: foundExpr.Pos(),
						End: foundExpr.End(),
						Message: fmt.Sprintf(
							"using %s from the outer scope but %s is defined inner at %s",
							types.ExprString(foundExpr),
							objInScope.Name(),
							pass.Fset.Position(objInScope.Pos()),
						),
						SuggestedFixes: []analysis.SuggestedFix{
							{
								Message: "use inner value of the same type",
								TextEdits: []analysis.TextEdit{
									{
										Pos:     foundExpr.Pos(),
										End:     foundExpr.End(),
										NewText: []byte(objInScope.Name()),
									},
								},
							},
						},
					},
				)

				return true
			}

			thisScope = thisScope.Parent()
		}

		return true
	})

	return nil, nil
}

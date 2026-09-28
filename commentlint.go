// Package commentlint limits the size of comments relative to the code they describe.
package commentlint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("commentlint", New)
}

type limits struct {
	MaxLines        int     `json:"max-lines"`
	Ratio           float64 `json:"ratio"`
	MinLines        int     `json:"min-lines"`
	ComplexityRatio float64 `json:"complexity-ratio"`
}

// declLimits has no complexity-ratio, so strict decoding rejects that key under decls.
type declLimits struct {
	MaxLines int     `json:"max-lines"`
	Ratio    float64 `json:"ratio"`
	MinLines int     `json:"min-lines"`
}

type settings struct {
	Funcs limits     `json:"funcs"`
	Decls declLimits `json:"decls"`
}

// New builds the plugin from the raw settings golangci-lint decodes from .golangci.yml.
func New(rawSettings any) (register.LinterPlugin, error) {
	s := settings{
		Funcs: limits{MaxLines: 3, Ratio: 0.15, MinLines: 1},
		Decls: declLimits{MaxLines: 2, Ratio: 1, MinLines: 1},
	}
	raw, err := json.Marshal(rawSettings)
	if err != nil {
		return nil, fmt.Errorf("commentlint: encoding settings: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&s); err != nil {
		return nil, fmt.Errorf("commentlint: decoding settings: %w", err)
	}
	for _, v := range []struct {
		key   string
		value float64
	}{
		{"funcs.max-lines", float64(s.Funcs.MaxLines)},
		{"funcs.ratio", s.Funcs.Ratio},
		{"funcs.min-lines", float64(s.Funcs.MinLines)},
		{"funcs.complexity-ratio", s.Funcs.ComplexityRatio},
		{"decls.max-lines", float64(s.Decls.MaxLines)},
		{"decls.ratio", s.Decls.Ratio},
		{"decls.min-lines", float64(s.Decls.MinLines)},
	} {
		if v.value < 0 {
			return nil, fmt.Errorf("commentlint: %s must not be negative, got %v", v.key, v.value)
		}
	}
	return &plugin{settings: s}, nil
}

type plugin struct {
	settings settings
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{{
		Name: "commentlint",
		Doc:  "reports comments that are larger than the code they describe allows",
		Run:  p.run,
	}}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

func (p *plugin) run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Doc == nil {
				continue
			}
			if lines, allowed := len(fn.Doc.List), p.settings.Funcs.MaxLines; allowed > 0 && lines > allowed {
				pass.Reportf(fn.Doc.Pos(), "func doc has %d comment lines, allowed %d (funcs.max-lines)", lines, allowed)
			}
		}
	}
	return nil, nil
}

package commentlint_test

import (
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	commentlint "github.com/prosayfer/comment-lint"
)

func run(t *testing.T, settings any, pkg string) {
	t.Helper()
	plugin, err := commentlint.New(settings)
	if err != nil {
		t.Fatal(err)
	}
	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, analysistest.TestData(), analyzers[0], pkg)
}

func TestFuncDocMaxLinesDefault(t *testing.T) {
	run(t, map[string]any{"funcs": map[string]any{"ratio": 0}}, "maxlines")
}

func TestFuncDocMaxLinesOverride(t *testing.T) {
	run(t, map[string]any{"funcs": map[string]any{"ratio": 0, "max-lines": 5}}, "maxlinesoverride")
}

func TestZeroDisablesMaxLines(t *testing.T) {
	run(t, map[string]any{"funcs": map[string]any{"ratio": 0, "max-lines": 0}}, "maxlineszero")
}

func TestInvalidSettings(t *testing.T) {
	tests := map[string]struct {
		settings map[string]any
		wantErr  string
	}{
		"unknown top-level key": {map[string]any{"func": map[string]any{}}, `"func"`},
		"unknown group key":     {map[string]any{"funcs": map[string]any{"max-line": 2}}, `"max-line"`},
		"complexity in decls":   {map[string]any{"decls": map[string]any{"complexity-ratio": 0.1}}, `"complexity-ratio"`},
		"negative max-lines":    {map[string]any{"funcs": map[string]any{"max-lines": -1}}, "funcs.max-lines"},
		"negative ratio":        {map[string]any{"decls": map[string]any{"ratio": -0.5}}, "decls.ratio"},
		"negative min-lines":    {map[string]any{"decls": map[string]any{"min-lines": -1}}, "decls.min-lines"},
		"negative complexity":   {map[string]any{"funcs": map[string]any{"complexity-ratio": -0.1}}, "funcs.complexity-ratio"},
		"string ratio":          {map[string]any{"funcs": map[string]any{"ratio": "high"}}, "funcs.ratio"},
		"fractional max-lines":  {map[string]any{"decls": map[string]any{"max-lines": 1.5}}, "decls.max-lines"},
		"group not a map":       {map[string]any{"funcs": 3}, "funcs"},
		"group in wrong case":   {map[string]any{"FUNCS": map[string]any{}}, `"FUNCS"`},
		"key in wrong case":     {map[string]any{"funcs": map[string]any{"Max-Lines": 5}}, `"funcs.Max-Lines"`},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := commentlint.New(tt.settings)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("New() error = %v, want it to mention %s", err, tt.wantErr)
			}
		})
	}
}

func TestNewAnalyzerRejectsNegative(t *testing.T) {
	s := commentlint.DefaultSettings()
	s.Funcs.ComplexityRatio = -0.1
	if _, err := commentlint.NewAnalyzer(s); err == nil || !strings.Contains(err.Error(), "funcs.complexity-ratio") {
		t.Fatalf("NewAnalyzer() error = %v, want it to mention funcs.complexity-ratio", err)
	}
}

func TestCommentLineCounting(t *testing.T) {
	run(t, map[string]any{"funcs": map[string]any{"ratio": 0}}, "counting")
}

func TestFuncDocFormulaDefaults(t *testing.T) {
	run(t, nil, "formula")
}

func TestFuncDocMinLinesOverride(t *testing.T) {
	run(t, map[string]any{"funcs": map[string]any{"min-lines": 2}}, "minlines")
}

func TestDeclDocsDefaults(t *testing.T) {
	run(t, nil, "decls")
}

func TestStructDocExcludesFieldComments(t *testing.T) {
	run(t, map[string]any{"decls": map[string]any{"max-lines": 0}}, "declsuncapped")
}

func TestExemptionsAndSkippedFiles(t *testing.T) {
	run(t, nil, "exempt")
}

func TestInBodyAndTrailingComments(t *testing.T) {
	run(t, nil, "inbody")
}

func TestFuncDocComplexityCeiling(t *testing.T) {
	run(t, map[string]any{"funcs": map[string]any{"complexity-ratio": 0.1}}, "complexity")
}

func TestFuncDocComplexityCeilingOffByDefault(t *testing.T) {
	run(t, nil, "complexityoff")
}

func TestComplexityCeilingSkipsInBodyAndTrailingComments(t *testing.T) {
	run(t, map[string]any{"funcs": map[string]any{"complexity-ratio": 0.1, "min-lines": 0}}, "complexityinbody")
}

func TestCRLFCountsLikeLF(t *testing.T) {
	run(t, nil, "crlf")
}

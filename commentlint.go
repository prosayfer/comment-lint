// Package commentlint limits the size of comments relative to the code they describe.
package commentlint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/scanner"
	"go/token"
	"math"
	"strings"

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
	funcDoc := rule{group: "funcs", limits: p.settings.Funcs}
	d := p.settings.Decls
	declDoc := rule{group: "decls", limits: limits{MaxLines: d.MaxLines, Ratio: d.Ratio, MinLines: d.MinLines}}
	for _, f := range pass.Files {
		tf := pass.Fset.File(f.Pos())
		src, err := pass.ReadFile(tf.Name())
		if err != nil {
			return nil, err
		}
		c := checker{pass: pass, code: codeLines(tf, src)}
		for _, decl := range f.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				c.check(decl.Doc, "func doc", funcDoc, decl)
			case *ast.GenDecl:
				if decl.Tok == token.IMPORT {
					continue
				}
				c.check(decl.Doc, "decl doc", declDoc, decl)
				for _, spec := range decl.Specs {
					c.checkSpec(spec, declDoc)
				}
			}
		}
	}
	return nil, nil
}

type rule struct {
	group string
	limits
}

// allowance returns the allowed comment lines and the setting that decided them, or "" when no limit applies.
func (r rule) allowance(code int) (int, string) {
	allowed, by := math.MaxInt, ""
	if r.Ratio > 0 {
		allowed, by = ceil(float64(code)*r.Ratio), "ratio"
	}
	if r.MinLines > 0 && r.MinLines >= allowed {
		allowed, by = r.MinLines, "min-lines"
	}
	if r.MaxLines > 0 && r.MaxLines <= allowed {
		allowed, by = r.MaxLines, "max-lines"
	}
	return allowed, by
}

// ceil tolerates float error, so that 20 × 0.15 rounds up to 3 rather than 4.
func ceil(x float64) int {
	return int(math.Ceil(x - 1e-9))
}

type checker struct {
	pass *analysis.Pass
	code lineSet
}

func (c *checker) check(g *ast.CommentGroup, kind string, r rule, subject ast.Node) {
	if g == nil {
		return
	}
	allowed, by := r.allowance(c.code.count(subject.Pos(), subject.End()))
	if lines := commentLines(g); by != "" && lines > allowed {
		c.pass.Reportf(g.Pos(), "%s has %d comment lines, allowed %d (%s.%s)", kind, lines, allowed, r.group, by)
	}
}

func (c *checker) checkSpec(spec ast.Spec, r rule) {
	switch spec := spec.(type) {
	case *ast.ValueSpec:
		c.check(spec.Doc, "decl doc", r, spec)
	case *ast.TypeSpec:
		c.check(spec.Doc, "decl doc", r, spec)
	}
	ast.Inspect(spec, func(n ast.Node) bool {
		if st, ok := n.(*ast.StructType); ok {
			for _, field := range st.Fields.List {
				c.check(field.Doc, "field comment", r, field)
			}
		}
		return true
	})
}

// lineSet holds the lines of a file that contain code, keyed by line number.
type lineSet struct {
	file  *token.File
	lines map[int]bool
}

func codeLines(tf *token.File, src []byte) lineSet {
	set := lineSet{file: tf, lines: map[int]bool{}}
	var s scanner.Scanner
	s.Init(tf, src, nil, 0)
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			return set
		}
		if tok == token.SEMICOLON && lit == "\n" {
			continue
		}
		end := pos
		if tok == token.STRING {
			end += token.Pos(len(lit) - 1)
		}
		for line := set.line(pos); line <= set.line(end); line++ {
			set.lines[line] = true
		}
	}
}

// line ignores //line directives, which would otherwise remap positions.
func (s lineSet) line(p token.Pos) int {
	return s.file.PositionFor(p, false).Line
}

func (s lineSet) count(from, to token.Pos) int {
	n := 0
	for line := s.line(from); line <= s.line(to); line++ {
		if s.lines[line] {
			n++
		}
	}
	return n
}

func commentLines(g *ast.CommentGroup) int {
	n := 0
	for _, c := range g.List {
		if text, ok := strings.CutPrefix(c.Text, "//"); ok {
			if !isDirective(text) && strings.TrimSpace(text) != "" {
				n++
			}
			continue
		}
		for _, line := range strings.Split(strings.TrimSuffix(strings.TrimPrefix(c.Text, "/*"), "*/"), "\n") {
			if strings.TrimSpace(line) != "" {
				n++
			}
		}
	}
	return n
}

func isDirective(text string) bool {
	for _, prefix := range []string{"go:", "nolint", "export ", "line ", " +build"} {
		if strings.HasPrefix(text, prefix) {
			return true
		}
	}
	return false
}

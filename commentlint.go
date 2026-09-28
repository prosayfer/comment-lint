// Package commentlint limits the size of comments relative to the code they describe.
package commentlint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/scanner"
	"go/token"
	"go/types"
	"math"
	"slices"
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
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
	funcs := rule{group: "funcs", limits: p.settings.Funcs}
	d := p.settings.Decls
	decls := rule{group: "decls", limits: limits{MaxLines: d.MaxLines, Ratio: d.Ratio, MinLines: d.MinLines}}
	for _, f := range pass.Files {
		tf := pass.Fset.File(f.Pos())
		if strings.HasSuffix(tf.Name(), "_test.go") || ast.IsGenerated(f) {
			continue
		}
		src, err := pass.ReadFile(tf.Name())
		if err != nil {
			return nil, err
		}
		c := checker{pass: pass, code: codeLines(tf, src), funcs: funcs, decls: decls}
		for _, decl := range f.Decls {
			c.checkDecl(decl)
		}
		for _, g := range f.Comments {
			c.checkInBody(f, g)
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
	pass         *analysis.Pass
	code         lineSet
	funcs, decls rule
}

func (c *checker) check(g *ast.CommentGroup, kind string, r rule, code int) {
	if g == nil {
		return
	}
	allowed, by := r.allowance(code)
	if lines := commentLines(g); by != "" && lines > allowed {
		c.pass.Reportf(g.Pos(), "%s has %d comment lines, allowed %d (%s.%s)", kind, lines, allowed, r.group, by)
	}
}

func (c *checker) checkDecl(decl ast.Decl) {
	switch decl := decl.(type) {
	case *ast.FuncDecl:
		c.check(decl.Doc, "func doc", c.funcs, c.code.count(decl))
	case *ast.GenDecl:
		if decl.Tok == token.IMPORT {
			return
		}
		if slices.ContainsFunc(decl.Specs, func(spec ast.Spec) bool { return !c.isContract(spec) }) {
			c.check(decl.Doc, "decl doc", c.decls, c.code.count(decl))
		}
		for _, spec := range decl.Specs {
			if !c.isContract(spec) {
				c.checkSpec(spec)
			}
		}
	}
}

func (c *checker) checkSpec(spec ast.Spec) {
	switch spec := spec.(type) {
	case *ast.ValueSpec:
		c.check(spec.Doc, "decl doc", c.decls, c.code.count(spec))
		c.check(spec.Comment, "trailing comment", c.decls, 1)
	case *ast.TypeSpec:
		c.check(spec.Doc, "decl doc", c.decls, c.code.count(spec))
		c.check(spec.Comment, "trailing comment", c.decls, 1)
	}
	ast.Inspect(spec, func(n ast.Node) bool {
		if st, ok := n.(*ast.StructType); ok {
			for _, field := range st.Fields.List {
				if _, isFunc := c.pass.TypesInfo.TypeOf(field.Type).Underlying().(*types.Signature); !isFunc {
					c.check(field.Doc, "field comment", c.decls, c.code.count(field))
					c.check(field.Comment, "trailing comment", c.decls, 1)
				}
			}
		}
		return true
	})
}

// checkInBody measures a comment inside a func body against the node it precedes within its innermost enclosing node.
func (c *checker) checkInBody(f *ast.File, g *ast.CommentGroup) {
	path, _ := astutil.PathEnclosingInterval(f, g.Pos(), g.End())
	if len(path) < 2 {
		return
	}
	if fn, ok := path[len(path)-2].(*ast.FuncDecl); !ok || fn.Body == nil || g.Pos() < fn.Body.Lbrace {
		return
	}
	if slices.ContainsFunc(path, func(n ast.Node) bool { _, ok := n.(*ast.InterfaceType); return ok }) {
		return
	}
	if c.code.startsBefore(g.Pos()) {
		c.check(g, "trailing comment", c.funcs, 1)
		return
	}
	var next ast.Node
	ast.Inspect(path[0], func(n ast.Node) bool {
		if n == path[0] {
			return true
		}
		if next == nil && n != nil && n.Pos() >= g.End() {
			next = n
		}
		return false
	})
	code := 1
	if next != nil {
		code = c.code.count(next)
	}
	c.check(g, "in-body comment", c.funcs, code)
}

// isContract reports whether spec declares an interface or func type, whose docs describe a contract and are exempt.
func (c *checker) isContract(spec ast.Spec) bool {
	ts, ok := spec.(*ast.TypeSpec)
	if !ok {
		return false
	}
	switch c.pass.TypesInfo.Defs[ts.Name].Type().Underlying().(type) {
	case *types.Interface, *types.Signature:
		return true
	}
	return false
}

// lineSet maps each line of a file that contains code to the position of its first code token.
type lineSet struct {
	file  *token.File
	lines map[int]token.Pos
}

func codeLines(tf *token.File, src []byte) lineSet {
	set := lineSet{file: tf, lines: map[int]token.Pos{}}
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
			if _, ok := set.lines[line]; !ok {
				set.lines[line] = pos
			}
		}
	}
}

// line ignores //line directives, which would otherwise remap positions.
func (s lineSet) line(p token.Pos) int {
	return s.file.PositionFor(p, false).Line
}

func (s lineSet) count(n ast.Node) int {
	count := 0
	for line := s.line(n.Pos()); line <= s.line(n.End()); line++ {
		if _, ok := s.lines[line]; ok {
			count++
		}
	}
	return count
}

// startsBefore reports whether code precedes p on its line.
func (s lineSet) startsBefore(p token.Pos) bool {
	first, ok := s.lines[s.line(p)]
	return ok && first < p
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

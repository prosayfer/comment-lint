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
	"github.com/uudashr/gocognit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
)

func init() {
	register.Plugin("commentlint", New)
}

type limits struct {
	MaxLines int     `json:"max-lines"`
	Ratio    float64 `json:"ratio"`
	MinLines int     `json:"min-lines"`
}

// funcLimits alone has complexity-ratio, so strict decoding rejects that key under decls.
type funcLimits struct {
	MaxLines        int     `json:"max-lines"`
	Ratio           float64 `json:"ratio"`
	MinLines        int     `json:"min-lines"`
	ComplexityRatio float64 `json:"complexity-ratio"`
}

type settings struct {
	Funcs funcLimits `json:"funcs"`
	Decls limits     `json:"decls"`
}

// New builds the plugin from the raw settings golangci-lint decodes from .golangci.yml.
func New(rawSettings any) (register.LinterPlugin, error) {
	s := settings{
		Funcs: funcLimits{MaxLines: 3, Ratio: 0.15, MinLines: 1},
		Decls: limits{MaxLines: 2, Ratio: 1, MinLines: 1},
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
	// These can't fail: raw just decoded into s, and s holds only numbers.
	var input, known map[string]any
	_ = json.Unmarshal(raw, &input)
	decoded, _ := json.Marshal(s)
	_ = json.Unmarshal(decoded, &known)
	if err := validate(input, known, ""); err != nil {
		return nil, err
	}
	return &plugin{settings: s}, nil
}

// validate rejects negative values and keys that differ from a known key only in case, which encoding/json accepts.
func validate(input, known map[string]any, path string) error {
	for key, value := range input {
		knownValue, ok := known[key]
		if !ok {
			return fmt.Errorf("commentlint: unknown key %q", path+key)
		}
		switch value := value.(type) {
		case map[string]any:
			if err := validate(value, knownValue.(map[string]any), path+key+"."); err != nil {
				return err
			}
		case float64:
			if value < 0 {
				return fmt.Errorf("commentlint: %s must not be negative, got %v", path+key, value)
			}
		}
	}
	return nil
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
	fl := p.settings.Funcs
	for _, f := range pass.Files {
		tf := pass.Fset.File(f.Pos())
		if strings.HasSuffix(tf.Name(), "_test.go") || ast.IsGenerated(f) {
			continue
		}
		src, err := pass.ReadFile(tf.Name())
		if err != nil {
			return nil, err
		}
		c := checker{
			pass:            pass,
			code:            scanCodeLines(tf, src),
			funcs:           rule{group: "funcs", limits: limits{MaxLines: fl.MaxLines, Ratio: fl.Ratio, MinLines: fl.MinLines}},
			decls:           rule{group: "decls", limits: p.settings.Decls},
			complexityRatio: fl.ComplexityRatio,
		}
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
	// complexityRatio caps the allowance by complexity, the cognitive complexity of the func; set only for func docs.
	complexityRatio float64
	complexity      int
}

// allowance returns the allowed comment lines and the setting that decided them, or "" when no limit applies.
func (r rule) allowance(code int) (int, string) {
	allowed, by := math.MaxInt, ""
	if r.Ratio > 0 {
		allowed, by = ceil(float64(code)*r.Ratio), "ratio"
	}
	if r.complexityRatio > 0 {
		if c := ceil(float64(r.complexity) * r.complexityRatio); c <= allowed {
			allowed, by = c, "complexity-ratio"
		}
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
	code         codeLines
	funcs, decls rule
	// complexityRatio applies only to func docs, so it is set on their rule together with the func's complexity.
	complexityRatio float64
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
		doc := c.funcs
		doc.complexityRatio, doc.complexity = c.complexityRatio, gocognit.Complexity(decl)
		c.check(decl.Doc, "func doc", doc, c.code.count(decl))
	case *ast.GenDecl:
		if decl.Tok == token.IMPORT {
			return
		}
		if !c.isContract(decl) {
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
		switch n := n.(type) {
		case *ast.FuncLit:
			// checkInBody owns everything inside a func literal.
			return false
		case *ast.StructType:
			for _, field := range n.Fields.List {
				if !c.isContract(field) {
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
	inBody := slices.ContainsFunc(path, func(n ast.Node) bool {
		var body *ast.BlockStmt
		switch n := n.(type) {
		case *ast.FuncDecl:
			body = n.Body
		case *ast.FuncLit:
			body = n.Body
		}
		return body != nil && body.Lbrace < g.Pos()
	})
	if !inBody || slices.ContainsFunc(path, func(n ast.Node) bool { _, ok := n.(*ast.InterfaceType); return ok }) {
		return
	}
	if first, ok := c.code.lines[c.code.line(g.Pos())]; ok && first < g.Pos() {
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
	if decl, ok := next.(*ast.DeclStmt); ok {
		next = decl.Decl
	}
	if next == nil {
		c.check(g, "in-body comment", c.funcs, 1)
	} else if !c.isContract(next) {
		c.check(g, "in-body comment", c.funcs, c.code.count(next))
	}
}

// isContract reports whether n declares only interface or func types, or is a func-typed field; their docs are exempt.
func (c *checker) isContract(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.GenDecl:
		return !slices.ContainsFunc(n.Specs, func(spec ast.Spec) bool { return !c.isContract(spec) })
	case *ast.TypeSpec:
		switch c.pass.TypesInfo.TypeOf(n.Type).Underlying().(type) {
		case *types.Interface, *types.Signature:
			return true
		}
	case *ast.Field:
		_, isFunc := c.pass.TypesInfo.TypeOf(n.Type).Underlying().(*types.Signature)
		return isFunc
	}
	return false
}

// codeLines maps each line of a file that contains code to the position of its first code token.
type codeLines struct {
	file  *token.File
	lines map[int]token.Pos
}

func scanCodeLines(tf *token.File, src []byte) codeLines {
	set := codeLines{file: tf, lines: map[int]token.Pos{}}
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
		// Counts newlines rather than bytes, because the scanner strips \r from raw strings.
		first := set.line(pos)
		for line := first; line <= first+strings.Count(lit, "\n"); line++ {
			if _, ok := set.lines[line]; !ok {
				set.lines[line] = pos
			}
		}
	}
}

// line ignores //line directives, which would otherwise remap positions.
func (s codeLines) line(p token.Pos) int {
	return s.file.PositionFor(p, false).Line
}

func (s codeLines) count(n ast.Node) int {
	count := 0
	for line := s.line(n.Pos()); line <= s.line(n.End()); line++ {
		if _, ok := s.lines[line]; ok {
			count++
		}
	}
	return count
}

func commentLines(g *ast.CommentGroup) int {
	n := 0
	for _, c := range g.List {
		if text, ok := strings.CutPrefix(c.Text, "//"); ok {
			directive := slices.ContainsFunc([]string{"go:", "nolint", "export ", "line ", " +build"}, func(prefix string) bool {
				return strings.HasPrefix(text, prefix)
			})
			if !directive && strings.TrimSpace(text) != "" {
				n++
			}
			continue
		}
		for _, line := range strings.Split(strings.TrimSuffix(strings.TrimPrefix(c.Text, "/*"), "*/"), "\n") {
			// A line of only `*` is decoration, as in /** … */ blocks.
			if strings.Trim(line, " \t*") != "" {
				n++
			}
		}
	}
	return n
}

// Command commentlint runs the commentlint analyzer with its default settings, without golangci-lint.
package main

import (
	"log"

	"golang.org/x/tools/go/analysis/singlechecker"

	commentlint "github.com/prosayfer/comment-lint"
)

func main() {
	analyzer, err := commentlint.NewAnalyzer(commentlint.DefaultSettings())
	if err != nil {
		log.Fatal(err)
	}
	singlechecker.Main(analyzer)
}

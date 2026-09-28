// Command commentlint runs the commentlint analyzer with its default settings, without golangci-lint.
package main

import (
	"log"

	"golang.org/x/tools/go/analysis/singlechecker"

	commentlint "github.com/prosayfer/comment-lint"
)

func main() {
	plugin, err := commentlint.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		log.Fatal(err)
	}
	singlechecker.Main(analyzers[0])
}

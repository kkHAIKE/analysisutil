package analysisutil_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestIsGeneratedFile(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		src  string
		want bool
	}{
		"true":  {"// Code generated by test; DO NOT EDIT.", true},
		"false": {"//Code generated by test; DO NOT EDIT.", false},
		"empty": {"", false},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			a := &analysis.Analyzer{
				Name: name + "Analyzer",
				Run: func(pass *analysis.Pass) (interface{}, error) {
					got := analysisutil.IsGeneratedFile(pass.Files[0])
					if tt.want != got {
						return nil, fmt.Errorf("want %v but got %v", tt.want, got)
					}
					return nil, nil
				},
			}
			path := filepath.Join(name, name+".go")
			dir := WriteFiles(t, map[string]string{
				path: fmt.Sprintf("%s\npackage %s", tt.src, name),
			})
			analysistest.Run(t, dir, a, name)
		})
	}
}

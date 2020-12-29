package mathjax

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"

	"github.com/stretchr/testify/assert"
)

type mathJaxTestCase struct {
	d   string // test description
	in  string // input markdown source
	out string // expected output html
}

func TestMathJax(t *testing.T) {

	tests := []mathJaxTestCase{
		{
			d:   "plain text",
			in:  "foo",
			out: `<p>foo</p>`,
		},
		{
			d:   "bold",
			in:  "**foo**",
			out: `<p><strong>foo</strong></p>`,
		},
		{
			d:   "math inline",
			in:  "$1+2$",
			out: `<p><span class="math inline">\(1+2\)</span></p>`,
		},
		{
			d:  "math display",
			in: "$$\n1+2\n$$",
			out: `<p><span class="math display">\[1+2
\]</span></p>`,
		},
		{
			// this input previously triggered a panic in block.go
			d:   "list-begin",
			in:  "*foo\n  ",
			out: "<p>*foo</p>",
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d: %s", i, tc.d), func(t *testing.T) {
			out, err := renderMarkdown(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.out, out)
		})
	}

}

func renderMarkdown(src string) (string, error) {
	var parserOptions []parser.Option
	var rendererOptions []renderer.Option

	extensions := []goldmark.Extender{
		MathJax,
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extensions...,
		),
		goldmark.WithParserOptions(
			parserOptions...,
		),
		goldmark.WithRendererOptions(
			rendererOptions...,
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

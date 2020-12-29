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
		// TODO: fix this bug.  This input triggers:
		//
		// 	panic: runtime error: index out of range [-1] [recovered]
		// 	panic: runtime error: index out of range [-1]

		// goroutine 24 [running]:
		// testing.tRunner.func1.1(0x1348dc0, 0xc000170fc0)
		// 	/usr/local/Cellar/go/1.15.2/libexec/src/testing/testing.go:1076 +0x30d
		// testing.tRunner.func1(0xc000083380)
		// 	/usr/local/Cellar/go/1.15.2/libexec/src/testing/testing.go:1079 +0x41a
		// panic(0x1348dc0, 0xc000170fc0)
		// 	/usr/local/Cellar/go/1.15.2/libexec/src/runtime/panic.go:969 +0x175
		// github.com/litao91/goldmark-mathjax.(*mathJaxBlockParser).Open(0x15b28e0, 0x13fdce0, 0xc000175580, 0x13fcee0, 0xc000192770, 0x13fd120, 0xc000192850, 0x1, 0xc000175600, 0x10)
		// 	/Users/foo/go/src/github.com/pcj/goldmark-mathjax/block.go:28 +0x195
		// github.com/yuin/goldmark/parser.(*parser).openBlocks(0xc000234a00, 0x13fdce0, 0xc000175580, 0x0, 0x13fcee0, 0xc000192770, 0x13fd120, 0xc000192850, 0x2)
		// 	/Users/foo/go/pkg/mod/github.com/yuin/goldmark@v1.2.1/parser/parser.go:938 +0x26c
		// github.com/yuin/goldmark/parser.(*parser).parseBlocks(0xc000234a00, 0x13fdce0, 0xc000175580, 0x13fcee0, 0xc000192770, 0x13fd120, 0xc000192850)
		// 	/Users/foo/go/pkg/mod/github.com/yuin/goldmark@v1.2.1/parser/parser.go:1082 +0x2df
		// github.com/yuin/goldmark/parser.(*parser).Parse(0xc000234a00, 0x13fcee0, 0xc000192770, 0x0, 0x0, 0x0, 0x8, 0x8)
		// 	/Users/foo/go/pkg/mod/github.com/yuin/goldmark@v1.2.1/parser/parser.go:848 +0x145
		// github.com/yuin/goldmark.(*markdown).Convert(0xc0000cef80, 0xc000173948, 0x7, 0x8, 0x13f73a0, 0xc00018ea80, 0x0, 0x0, 0x0, 0x149bd20, ...)
		// 	/Users/foo/go/pkg/mod/github.com/yuin/goldmark@v1.2.1/markdown.go:116 +0x96
		// github.com/litao91/goldmark-mathjax.renderMarkdown(0x13684d4, 0x7, 0x133be381, 0x958c133be381, 0x100000001, 0xc00003e748)
		// 	/Users/foo/go/src/github.com/pcj/goldmark-mathjax/mathjax_test.go:87 +0x1d4
		// github.com/litao91/goldmark-mathjax.TestMathJax.func1(0xc000083380)
		// 	/Users/foo/go/src/github.com/pcj/goldmark-mathjax/mathjax_test.go:56 +0x4e
		// testing.tRunner(0xc000083380, 0xc000220dd0)
		// 	/usr/local/Cellar/go/1.15.2/libexec/src/testing/testing.go:1127 +0xef
		// created by testing.(*T).Run
		// 	/usr/local/Cellar/go/1.15.2/libexec/src/testing/testing.go:1178 +0x386
		// FAIL	github.com/litao91/goldmark-mathjax	0.261s
		{
			d:   "list-bug",
			in:  "*foo\n  ",
			out: "<p>*foo</p>\n",
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

package mathjax

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type mathjax struct {
}

var MathJax = &mathjax{}

func (*mathjax) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(NewMathJaxBlockParser(), 701),
	))
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewInlineMathParser(), 501),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewMathBlockRenderer(), 501),
		util.Prioritized(NewInlineMathRenderer(), 502),
	))
}

package mathjax

import (
	"bytes"
	"fmt"
	"strings"

	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type MathTexBlockRenderer struct {
	renderer *TexRenderer
}

func (r *MathTexBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindMathBlock, r.renderMathBlock)
}

func (r *MathTexBlockRenderer) writeLines(w *bytes.Buffer, source []byte, n gast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		w.Write(line.Value(source))
	}
}

func (r *MathTexBlockRenderer) renderMathBlock(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*MathBlock)
	if entering {
		_, _ = w.WriteString(`<div class="latex-svg display" style="vertical-align:middle;">`)
		var buf bytes.Buffer
		r.writeLines(&buf, source, n)
		str := buf.String()
		fmt.Println("===" + str)
		var svg []byte
		if strings.Contains(str, "tikzpicture") {
			svg = r.renderer.Run(str)
		} else {
			svg = r.renderer.Run(`\[` + str + `\]`)
		}
		_, _ = w.WriteString(string(svg))
	} else {
		_, _ = w.WriteString(`</div>` + "\n")
	}
	return gast.WalkContinue, nil
}

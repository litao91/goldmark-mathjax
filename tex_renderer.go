package mathjax

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const common = `
{{.Doc}}
`

const displayInlineFormula = `\( {{.Doc}} \)`

const displayBlockFormula = `\[
{{.Doc}}
\]`

const tmpl = `
\documentclass[11pt]{article}
\usepackage[paperwidth=180in,paperheight=180in]{geometry}
\batchmode
\usepackage[utf8]{inputenc}
\usepackage{amsmath}
\usepackage{amssymb}
\usepackage{stmaryrd}
\usepackage[verbose]{newunicodechar}
\pagestyle{empty}
\setlength{\topskip}{0pt}
\setlength{\parindent}{0pt}
\setlength{\abovedisplayskip}{0pt}
\setlength{\belowdisplayskip}{0pt}

\begin{document}
{{.Doc}}
\end{document}
`

type TexRenderer struct {
	texPath         string
	docTemplate     *template.Template
	inlineFormulaImpl      *template.Template
	commonBlockTmpl *template.Template
	blockFormulaTmpl     *template.Template
	tmpDir          string
}

func NewDefaultTexRenderer() *TexRenderer {
	docTmpl, err := template.New("text").Parse(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	inlineTmpl, err := template.New("text").Parse(displayInlineFormula)
	if err != nil {
		fmt.Println(err)
	}
	displayBlockFormulaTmpl, err := template.New("text").Parse(displayBlockFormula)
	if err != nil {
		fmt.Println(err)
	}
	commonTmpl, err := template.New("text").Parse(common)
	if err != nil {
		fmt.Println(err)
	}

	var wd, _ = os.Getwd()
	var texPath = os.Getenv("TEX_PATH")

	var tmpDir = wd + "/tmp/"

	var defaultRenderer = &TexRenderer{
		texPath:         texPath,
		docTemplate:     docTmpl,
		inlineFormulaImpl:      inlineTmpl,
		commonBlockTmpl: commonTmpl,
		blockFormulaTmpl:     displayBlockFormulaTmpl,
		tmpDir:          tmpDir,
	}
	return defaultRenderer
}

type DocT struct {
	Doc string
}

func (r *TexRenderer) RunInline(formula string) []byte {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	r.inlineFormulaImpl.Execute(writer, &DocT{
		Doc: formula,
	})
	writer.Flush()

	var docBuf bytes.Buffer
	docWriter := bufio.NewWriter(&docBuf)
	r.docTemplate.Execute(docWriter, &DocT{
		Doc: buf.String(),
	})
	docWriter.Flush()
	return r.runRaw(docBuf.String())
}

func (r *TexRenderer) Run(formula string) []byte {
	var bf bytes.Buffer
	writer := bufio.NewWriter(&bf)
	var tmpl *template.Template
	if strings.Contains(formula, `\begin{`) {
		tmpl = r.commonBlockTmpl
	} else {
		tmpl = r.blockFormulaTmpl
	}
	tmpl.Execute(writer, &DocT{
		Doc: formula,
	})
	writer.Flush()

	var docBf bytes.Buffer
	docWriter := bufio.NewWriter(&docBf)
	r.docTemplate.Execute(docWriter, &DocT{bf.String()})
	docWriter.Flush()
	return r.runRaw(docBf.String())
}

func (r *TexRenderer) runRaw(formula string) []byte {
	f, err := ioutil.TempFile(r.tmpDir, "doc")
	if err != nil {
		log.Fatalf("%v", err)
	}
	f.WriteString(formula)
	f.Sync()
	f.Close()
	r.runLatex(f.Name())
	r.runDvi2Svg(f.Name())
	svgf, err := os.Open(f.Name() + ".svg")
	if err != nil {
		return nil
	}
	svg, err := ioutil.ReadAll(svgf)
	if err != nil {
		return nil
	}
	return svg
}

func (r *TexRenderer) runDvi2Svg(fname string) {
	// fmt.Println([]string{fmt.Sprintf("%sdvisvgm", r.texPath), fmt.Sprintf("%s.dvi", fname), "-o", fmt.Sprintf("%s.svg", fname), "-n", "--exact", "-v0", "--relative", "--zoom=1.2546875"})

	cmd := exec.Command(fmt.Sprintf("%sdvisvgm", r.texPath), fmt.Sprintf("%s.dvi", fname), "-o", fmt.Sprintf("%s.svg", fname), "-n", "--exact", "-v0", "--relative", "--zoom=1.2546875")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("dvi2svg cmd.Run() failed with %s\n", err)
	}
	// outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	// fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}

func (r *TexRenderer) runLatex(fname string) {
	cmd := exec.Command(fmt.Sprintf("%slatex", r.texPath), "-output-directory", r.tmpDir, fname)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("latex cmd.Run() failed with %s\n", err)
	}
	// outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	// fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}

func (r *TexRenderer) runPdfLatex(fname string) {
	cmd := exec.Command(fmt.Sprintf("%spdflatex", r.texPath), "-output-directory", r.tmpDir, fname)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("pdflatex cmd.Run() failed with %s\n", err)
	}
}

func (r *TexRenderer) runPdf2Svg(fname string) {
	cmd := exec.Command("pdf2svg", fmt.Sprintf("%s.pdf", fname), fmt.Sprintf("%s.svg", fname))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("pdf2svg cmd.Run() failed with %s\n", err)
	}
}

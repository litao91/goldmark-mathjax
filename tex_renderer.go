package mathjax

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

const tmpl = `
\documentclass[11pt]{article}
\usepackage[paperwidth=180in,paperheight=180in]{geometry}
\batchmode
\usepackage[utf8]{inputenc}
\usepackage{amsmath}
\usepackage{amssymb}
\usepackage{stmaryrd}
\newcommand{\R}{\mathbb{R}}

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
	texPath     string
	docTemplate *template.Template
	tmpDir      string
}

func NewDefaultTexRenderer() *TexRenderer {
	var t, err = template.New("html").Parse(tmpl)
	if err != nil {
		panic(err)
	}

	var wd, _ = os.Getwd()
	var texPath = os.Getenv("TEX_PATH")

	var tmpDir = wd + "/tmp/"

	var defaultRenderer = &TexRenderer{
		texPath:     texPath,
		docTemplate: t,
		tmpDir:      tmpDir,
	}
	return defaultRenderer
}

func (r *TexRenderer) Run(formula string) []byte {
	f, err := ioutil.TempFile(r.tmpDir, "doc")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer os.Remove(f.Name())
	r.docTemplate.Execute(f, struct {
		Doc string
	}{
		Doc: formula,
	})
	f.Sync()
	f.Close()
	r.runLatex(f.Name())
	r.runDvi2Svg(f.Name())
	svgf, err := os.Open(f.Name() + ".svg")
	if err != nil {
		return nil
	}
	svg, err := ioutil.ReadAll(svgf)
	if err != nil  {
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
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}

func (r *TexRenderer) runLatex(fname string) {
	cmd := exec.Command(fmt.Sprintf("%slatex", r.texPath), "-output-directory", r.tmpDir, fname)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}

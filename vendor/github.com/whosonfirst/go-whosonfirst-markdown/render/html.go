package render

import (
	"bytes"
	"github.com/whosonfirst/go-whosonfirst-markdown"
	"github.com/whosonfirst/go-whosonfirst-markdown/jekyll"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"io"
	"log"
)

type HTMLOptions struct {
	Mode      string
	Input     string
	Output    string
	Header    string
	Footer    string
	Templates *template.Template
}

func DefaultHTMLOptions() *HTMLOptions {

	opts := HTMLOptions{
		Mode:      "files",
		Input:     "index.md",
		Output:    "index.html",
		Header:    "",
		Footer:    "",
		Templates: nil,
	}

	return &opts
}

type nopCloser struct {
	io.Reader
}

type WOFRenderer struct {
	bf          *blackfriday.HTMLRenderer
	frontmatter *jekyll.FrontMatter
	header      string
	footer      string
	templates   *template.Template
}

func (r *WOFRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {

	switch node.Type {

	case blackfriday.Image:
		return r.bf.RenderNode(w, node, entering)
	default:
		return r.bf.RenderNode(w, node, entering)
	}
}

func (r *WOFRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {

	if r.templates == nil || r.header == "" {
		r.bf.RenderHeader(w, ast)
		return
	}

	t := r.templates.Lookup(r.header)

	if t == nil {
		log.Printf("Invalid or missing template '%s'\n", r.header)
		return
	}

	err := t.Execute(w, r.frontmatter)

	if err != nil {
		log.Println(err)
	}
}

func (r *WOFRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {

	if r.templates == nil || r.footer == "" {
		r.bf.RenderFooter(w, ast)
		return
	}

	t := r.templates.Lookup(r.footer)

	if t == nil {
		log.Printf("Invalid or missing template '%s'\n", r.footer)
		return
	}

	err := t.Execute(w, r.frontmatter)

	if err != nil {
		log.Println(err)
	}
}

func (nopCloser) Close() error { return nil }

func RenderHTML(d *markdown.Document, opts *HTMLOptions) (io.ReadCloser, error) {

	flags := blackfriday.CommonHTMLFlags
	flags |= blackfriday.CompletePage
	flags |= blackfriday.UseXHTML

	params := blackfriday.HTMLRendererParameters{
		Flags: flags,
	}

	renderer := blackfriday.NewHTMLRenderer(params)

	r := WOFRenderer{
		bf:          renderer,
		frontmatter: d.FrontMatter,
		header:      opts.Header,
		footer:      opts.Footer,
		templates:   opts.Templates,
	}

	unsafe := blackfriday.Run(d.Body.Bytes(), blackfriday.WithRenderer(&r))

	// safe := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	html := bytes.NewReader(unsafe)
	return nopCloser{html}, nil

}

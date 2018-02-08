package main

import (
	"context"
	"errors"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-markdown"
	"github.com/whosonfirst/go-whosonfirst-markdown/flags"
	"github.com/whosonfirst/go-whosonfirst-markdown/parser"
	"github.com/whosonfirst/go-whosonfirst-markdown/render"
	"github.com/whosonfirst/go-whosonfirst-markdown/writer"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func RenderDirectory(ctx context.Context, dir string, opts *render.HTMLOptions) error {

	cb := func(path string, info os.FileInfo) error {

		select {
		case <-ctx.Done():
			return nil
		default:

			if info.IsDir() {
				return nil
			}

			return RenderPathWithRoot(ctx, path, dir, opts)
		}
	}

	c := crawl.NewCrawler(dir)
	return c.Crawl(cb)
}

func RenderPath(ctx context.Context, path string, opts *render.HTMLOptions) error {
	root := filepath.Dir(path)
	return RenderPathWithRoot(ctx, path, root, opts)
}

func RenderPathWithRoot(ctx context.Context, path string, root string, opts *render.HTMLOptions) error {

	select {

	case <-ctx.Done():
		return nil
	default:

		abs_path, err := filepath.Abs(path)

		if err != nil {
			return err
		}

		fname := filepath.Base(abs_path)

		if fname != opts.Input {
			return nil
		}

		parse_opts := parser.DefaultParseOptions()
		fm, body, err := parser.ParseFile(abs_path, parse_opts)

		if err != nil {
			return err
		}

		doc, err := markdown.NewDocument(fm, body)
		html, err := render.RenderHTML(doc, opts)

		if err != nil {
			return err
		}

		wr := ctx.Value("writer").(writer.Writer)

		if wr == nil {
			return errors.New("Can't load writer from context")
		}

		// I don't love that all this logic is here but I am not
		// sure where else to put it... (20180109/thisisaaronland)

		out_path := fm.Permalink

		if out_path == "" {
			abs_root := filepath.Dir(abs_path)
			out_path = filepath.Join(abs_root, opts.Output)
		}

		if strings.HasSuffix(out_path, "/") {
			out_path = filepath.Join(out_path, opts.Output)
		}

		if root != "" {
			out_path = strings.Replace(out_path, root, "", -1)
		}

		return wr.Write(out_path, html)
	}
}

func Render(ctx context.Context, path string, opts *render.HTMLOptions) error {

	select {
	case <-ctx.Done():
		return nil
	default:

		abs_path, err := filepath.Abs(path)

		if err != nil {
			return err
		}

		switch opts.Mode {

		case "files":
			return RenderPath(ctx, abs_path, opts)
		case "directory":
			return RenderDirectory(ctx, abs_path, opts)
		default:
			return errors.New("Unknown or invalid mode")
		}
	}

}

func main() {

	var mode = flag.String("mode", "files", "Valid modes are: files, directory")
	var input = flag.String("input", "index.md", "What you expect the input Markdown file to be called")
	var output = flag.String("output", "index.html", "What you expect the output HTML file to be called")
	var header = flag.String("header", "", "The name of the (Go) template to use as a custom header")
	var footer = flag.String("footer", "", "The name of the (Go) template to use as a custom footer")

	var templates flags.HTMLTemplateFlags
	flag.Var(&templates, "templates", "One or more directories containing (Go) templates to parse")

	var writers flags.WriterFlags
	flag.Var(&writers, "writer", "One or more writer to output rendered Markdown to. Valid writers are: fs=PATH; null; stdout")

	flag.Parse()

	wr, err := writers.ToWriter()

	if err != nil {
		log.Fatal(err)
	}

	t, err := templates.Parse()

	if err != nil {
		log.Fatal(err)
	}

	opts := render.DefaultHTMLOptions()
	opts.Mode = *mode
	opts.Input = *input
	opts.Output = *output
	opts.Header = *header
	opts.Footer = *footer
	opts.Templates = t

	ctx := context.Background()
	ctx = context.WithValue(ctx, "writer", wr)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, path := range flag.Args() {

		err := Render(ctx, path, opts)

		if err != nil {
			log.Println(err)
			cancel()
			break
		}
	}
}

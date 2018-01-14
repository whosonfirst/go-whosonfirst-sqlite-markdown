package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-markdown"
	"github.com/whosonfirst/go-whosonfirst-markdown/jekyll"
	"github.com/whosonfirst/go-whosonfirst-markdown/parser"
	"github.com/whosonfirst/go-whosonfirst-markdown/render"
	"github.com/whosonfirst/go-whosonfirst-markdown/utils"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"text/template"
	"time"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func RenderDirectory(ctx context.Context, path string, opts *render.HTMLOptions) error {

	lookup := make(map[string]*jekyll.FrontMatter)
	dates := make([]string, 0)

	mu := new(sync.Mutex)

	cb := func(path string, info os.FileInfo) error {

		select {
		case <-ctx.Done():
			log.Println("DONE")
			return nil
		default:

			if info.IsDir() {
				return nil // PLEASE MAKE ME WORK... (20180113/thisisaaronland)
			}

			fm, err := RenderPath(ctx, path, opts)

			if err != nil {
				log.Println("OOPS", path, err)
				return err
			}

			if fm == nil {
				return nil
			}

			t, err := time.Parse("2006-01-02", fm.Date)

			if err != nil {
				return err
			}

			ymd := t.Format("20060102")

			mu.Lock()
			dates = append(dates, ymd)
			lookup[ymd] = fm
			mu.Unlock()

			return nil
		}
	}

	c := crawl.NewCrawler(path)
	c.CrawlDirectories = true

	err := c.Crawl(cb)

	if err != nil {
		log.Println("FAILED TO CRAWL", err)
		return nil
	}

	posts := make([]*jekyll.FrontMatter, 0)
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	for _, ymd := range dates {
		posts = append(posts, lookup[ymd])
	}

	if len(posts) == 0 {
		return nil
	}

	return RenderPosts(ctx, path, posts, opts)
}

func RenderPosts(ctx context.Context, root string, posts []*jekyll.FrontMatter, opts *render.HTMLOptions) error {

	select {
	case <-ctx.Done():
		return nil
	default:

		tm := `{{ range $fm := .Posts }}
### [{{ $fm.Title }}]({{ $fm.Permalink }})

> {{ $fm.Excerpt }}

> _{{ $fm.Date }}_
	    {{ end }}`

		t, err := template.New("index").Parse(tm)

		if err != nil {
			return err
		}

		type Data struct {
			Posts []*jekyll.FrontMatter
		}

		d := Data{
			Posts: posts,
		}

		var b bytes.Buffer
		wr := bufio.NewWriter(&b)

		err = t.Execute(wr, d)

		if err != nil {
			return err
		}

		r := bytes.NewReader(b.Bytes())
		fh := nopCloser{r}

		p_opts := parser.DefaultParseOptions()
		fm, buf, err := parser.Parse(fh, p_opts)

		if err != nil {
			return err
		}

		doc, err := markdown.NewDocument(fm, buf)

		if err != nil {
			return err
		}

		html, err := render.RenderHTML(doc, opts)

		if err != nil {
			return err
		}

		return utils.WriteHTML(html, root, opts)
	}
}

func RenderPath(ctx context.Context, path string, opts *render.HTMLOptions) (*jekyll.FrontMatter, error) {

	select {

	case <-ctx.Done():
		return nil, nil
	default:

		abs_path, err := filepath.Abs(path)

		if err != nil {
			return nil, err
		}

		parse_opts := parser.DefaultParseOptions()
		parse_opts.Body = false

		fm, _, err := parser.ParseFile(abs_path, parse_opts)

		if err != nil {
			return nil, err
		}

		return fm, nil
	}
}

func Render(ctx context.Context, path string, opts *render.HTMLOptions) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		return RenderDirectory(ctx, path, opts)
	}

}

func main() {

	var input = flag.String("input", "index.md", "What you expect the input Markdown file to be called")
	var output = flag.String("output", "index.html", "What you expect the output HTML file to be called")
	var header = flag.String("header", "", "The path to a custom (Go) template to use as header for your HTML output")
	var footer = flag.String("footer", "", "The path to a custom (Go) template to use as a footer for your HTML output")

	flag.Parse()

	opts := render.DefaultHTMLOptions()
	opts.Input = *input
	opts.Output = *output

	if *header != "" {

		t, err := utils.LoadTemplate(*header, "header")

		if err != nil {
			log.Fatal(err)
		}

		opts.Header = t
	}

	if *footer != "" {

		t, err := utils.LoadTemplate(*footer, "footer")

		if err != nil {
			log.Fatal(err)
		}

		opts.Footer = t
	}

	ctx, cancel := context.WithCancel(context.Background())
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

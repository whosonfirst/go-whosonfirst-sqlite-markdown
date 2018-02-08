package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-markdown/flags"
	"github.com/whosonfirst/go-whosonfirst-markdown/jekyll"
	"github.com/whosonfirst/go-whosonfirst-markdown/parser"
	"github.com/whosonfirst/go-whosonfirst-markdown/render"
	"github.com/whosonfirst/go-whosonfirst-markdown/writer"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func Render(ctx context.Context, path string, opts *render.FeedOptions) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		return RenderDirectory(ctx, path, opts)
	}

}

func RenderDirectory(ctx context.Context, dir string, opts *render.FeedOptions) error {

	posts, err := GatherPosts(ctx, dir, opts)

	if err != nil {
		return err
	}

	if len(posts) == 0 {
		return nil
	}

	return RenderPosts(ctx, dir, posts, opts)
}

// THIS IS A BAD NAME - ALSO SHOULD BE SHARED CODE...
// (20180130/thisisaaronland)

func RenderPath(ctx context.Context, path string, opts *render.FeedOptions) (*jekyll.FrontMatter, error) {

	select {

	case <-ctx.Done():
		return nil, nil
	default:

		abs_path, err := filepath.Abs(path)

		if err != nil {
			log.Printf("FAILED to render path %s, because %s\n", path, err)
			return nil, err
		}

		if filepath.Base(abs_path) != opts.Input {
			return nil, nil
		}

		parse_opts := parser.DefaultParseOptions()
		parse_opts.Body = false

		fm, _, err := parser.ParseFile(abs_path, parse_opts)

		if err != nil {
			log.Printf("FAILED to parse %s, because %s\n", path, err)
			return nil, err
		}

		return fm, nil
	}
}

func GatherPosts(ctx context.Context, root string, opts *render.FeedOptions) ([]*jekyll.FrontMatter, error) {

	mu := new(sync.Mutex)

	lookup := make(map[string]*jekyll.FrontMatter)
	dates := make([]string, 0)

	cb := func(path string, info os.FileInfo) error {

		select {
		case <-ctx.Done():
			return nil
		default:

			if info.IsDir() {
				return nil
			}

			abs_path, err := filepath.Abs(path)

			if err != nil {
				return err
			}

			if filepath.Base(abs_path) != opts.Input {
				return nil
			}

			fm, err := RenderPath(ctx, path, opts)

			if err != nil {
				return err
			}

			if fm == nil {
				return nil
			}

			mu.Lock()
			ymd := fm.Date.Format("20060102")
			dates = append(dates, ymd)
			lookup[ymd] = fm
			mu.Unlock()
		}

		return nil
	}

	c := crawl.NewCrawler(root)
	c.CrawlDirectories = true

	err := c.Crawl(cb)

	if err != nil {
		return nil, err
	}

	posts := make([]*jekyll.FrontMatter, 0)

	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	for _, ymd := range dates {
		posts = append(posts, lookup[ymd])

		if len(posts) == opts.Items {
			break
		}
	}

	return posts, nil
}

func RenderPosts(ctx context.Context, root string, posts []*jekyll.FrontMatter, opts *render.FeedOptions) error {

	select {
	case <-ctx.Done():
		return nil
	default:

		type Data struct {
			Posts     []*jekyll.FrontMatter
			BuildDate time.Time
		}

		now := time.Now()

		d := Data{
			Posts:     posts,
			BuildDate: now,
		}

		// PLEASE REPLACE ALL OF THIS WILL A GENERIC utils.WriteTemplate
		// THAT WRAPS ALL THE atomicfile STUFF AND WRITES STRAIGHT TO A
		// FILEHANDLE RATHER THAN BYTES THEN... (20180130/thisisaaronland)

		var b bytes.Buffer
		wr := bufio.NewWriter(&b)

		t_name := fmt.Sprintf("feed_%s", opts.Format)
		t := opts.Templates.Lookup(t_name)

		if t == nil {
			return errors.New(fmt.Sprintf("Invalid or missing template '%s'", t_name))
		}

		err := t.Execute(wr, d)

		if err != nil {
			return err
		}

		wr.Flush()

		r := bytes.NewReader(b.Bytes())
		fh := nopCloser{r}

		w := ctx.Value("writer").(writer.Writer)

		if w == nil {
			return errors.New("Can't load writer from context")
		}

		out_path := filepath.Join(root, opts.Output)
		return w.Write(out_path, fh)
	}
}

func main() {

	var input = flag.String("input", "index.md", "What you expect the input Markdown file to be called")
	var output = flag.String("output", "", "The filename of your feed. If empty default to the value of -format + \".xml\"")

	var format = flag.String("format", "rss_20", "Valid options are: atom_10, rss_20")
	var items = flag.Int("items", 10, "The number of items to include in your feed")

	var templates flags.FeedTemplateFlags
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

	if *output == "" {
		*output = fmt.Sprintf("%s.xml", *format)
	}

	opts := render.DefaultFeedOptions()
	opts.Input = *input
	opts.Output = *output
	opts.Format = *format
	opts.Items = *items
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

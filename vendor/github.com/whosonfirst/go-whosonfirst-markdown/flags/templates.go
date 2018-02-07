package flags

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	html_template "html/template"
	_ "log"
	"os"
	"path/filepath"
	"sync"
	text_template "text/template"
)

type HTMLTemplateFlags []string

func (t *HTMLTemplateFlags) String() string {
	return fmt.Sprintf("%v", *t)
}

func (t *HTMLTemplateFlags) Set(root string) error {

	mu := new(sync.Mutex)

	cb := func(path string, info os.FileInfo) error {

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)

		if ext != ".html" {
			return nil
		}

		mu.Lock()
		*t = append(*t, path)
		mu.Unlock()

		return nil
	}

	c := crawl.NewCrawler(root)
	return c.Crawl(cb)
}

// Maybe move this in to a different package
// (20180129/thisisaaronland)

func (t *HTMLTemplateFlags) Parse() (*html_template.Template, error) {

	if len(*t) == 0 {
		return nil, nil
	}

	// https://play.golang.org/p/V94BPN0uKD

	var fns = html_template.FuncMap{
		"plus1": func(x int) int {
			return x + 1
		},
	}

	// we need something attach Funcs() to before we call
	// ParseFiles() and if there's another way I don't know
	// what it is...

	return html_template.New("debug").Funcs(fns).ParseFiles(*t...)
}

type FeedTemplateFlags []string

func (t *FeedTemplateFlags) String() string {
	return fmt.Sprintf("%v", *t)
}

func (t *FeedTemplateFlags) Set(root string) error {

	mu := new(sync.Mutex)

	cb := func(path string, info os.FileInfo) error {

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)

		if ext != ".xml" {
			return nil
		}

		mu.Lock()
		*t = append(*t, path)
		mu.Unlock()

		return nil
	}

	c := crawl.NewCrawler(root)
	return c.Crawl(cb)
}

func (t *FeedTemplateFlags) Parse() (*text_template.Template, error) {

	if len(*t) == 0 {
		return nil, nil
	}

	// PLEASE RECONCILE THIS WITH ABOVE...

	var fns = text_template.FuncMap{
		"plus1": func(x int) int {
			return x + 1
		},
	}

	return text_template.New("debug").Funcs(fns).ParseFiles(*t...)
}

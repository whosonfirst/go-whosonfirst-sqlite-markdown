package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/djherbis/times"
	"github.com/whosonfirst/go-whosonfirst-markdown"
	"github.com/whosonfirst/go-whosonfirst-markdown/jekyll"
	"io"
	_ "log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type ParseOptions struct {
	FrontMatter bool
	Body        bool
}

func DefaultParseOptions() *ParseOptions {

	opts := ParseOptions{
		FrontMatter: true,
		Body:        true,
	}

	return &opts
}

func ParseFile(path string, opts *ParseOptions) (*jekyll.FrontMatter, *markdown.Body, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return nil, nil, err
	}

	fh, err := os.Open(abs_path)

	if err != nil {
		return nil, nil, err
	}

	defer fh.Close()

	fm, body, err := Parse(fh, opts)

	if err != nil {
		return nil, nil, err
	}

	if fm.Date == "" {

		re, err := regexp.Compile(`.*\/(\d{4})\/(\d{2})\/(\d{2})\/.*`)

		if err != nil {
			return nil, nil, err
		}

		m := re.FindAllStringSubmatch(abs_path, 1)

		if len(m) == 1 {

			yyyy := m[0][1]
			mm := m[0][2]
			dd := m[0][3]

			dt := fmt.Sprintf("%s-%s-%s", yyyy, mm, dd)
			fm.Date = dt
		} else {

			info, err := times.Stat(abs_path)

			if err != nil {
				return nil, nil, err
			}

			var t time.Time

			if info.HasBirthTime() {
				t = info.BirthTime()
			} else {
				t = info.ChangeTime() // not an awesome solution but what else can we do...
			}

			dt := t.Format("2006-01-02")
			fm.Date = dt
		}
	}

	if fm.Permalink == "" {

		// THIS IS ALL DEPRECATED

		/*
			fname := filepath.Base(abs_path)

			if fname != opts.Input {
				return nil, nil
			}

			root := filepath.Dir(abs_path)

			parts := strings.Split(root, "/")
			count := len(parts)

			yyyy := parts[(count-1)-3]
			mm := parts[(count-1)-2]
			dd := parts[(count-1)-1]
			post := parts[(count - 1)]

			uri := fmt.Sprintf("/blog/%s/%s/%s/%s/", yyyy, mm, dd, post)
		*/

		// END OF THIS IS ALL DEPRECATED

		// fm.Permalink = uri // FIX ME
	}

	return fm, body, nil
}

func Parse(md io.ReadCloser, opts *ParseOptions) (*jekyll.FrontMatter, *markdown.Body, error) {

	fm := jekyll.EmptyFrontMatter()

	var b bytes.Buffer
	wr := bufio.NewWriter(&b)

	scanner := bufio.NewScanner(md)

	lineno := 0
	is_jekyll := false

	for scanner.Scan() {

		lineno += 1

		txt := scanner.Text()
		ln := strings.Trim(txt, " ")

		if lineno == 1 && txt == "---" {
			is_jekyll = true
			continue
		}

		if is_jekyll && txt == "---" {
			is_jekyll = false
			continue
		}

		// this stuff should probably be moved in to jekyll/jekyll.go
		// (20181011/thisisaaronland)

		if is_jekyll {

			if opts.FrontMatter {

				kv := strings.Split(ln, ":")
				key := strings.Trim(kv[0], " ")
				value := strings.Trim(kv[1], " ")

				// log.Println("FRONT MATTER", ln)

				switch key {
				case "authors":
					fm.Authors = string2list(value)
				case "category":
					fm.Category = string2string(value)
				case "date":
					fm.Date = string2string(value)
				case "excerpt":
					fm.Excerpt = string2string(value)
				case "image":
					fm.Image = string2string(value)
				case "layout":
					fm.Layout = string2string(value)
				case "permalink":
					fm.Permalink = string2string(value)
				case "published":
					fm.Published = string2bool(value)
				case "tag":
					fm.Tags = string2list(value)
				case "tags":
					fm.Tags = string2list(value)
				case "title":
					fm.Title = string2string(value)
				default:
					// pass
				}
			}
			continue
		}

		if opts.Body {
			wr.WriteString(txt + "\n")
		}
	}

	wr.Flush()
	body := markdown.Body{&b}

	return fm, &body, nil
}

func string2string(s string) string {
	s = strings.TrimLeft(s, "\"")
	s = strings.TrimRight(s, "\"")
	return s
}

func string2list(s string) []string {
	s = strings.TrimLeft(s, "[")
	s = strings.TrimRight(s, "]")

	l := make([]string, 0)

	for _, str := range strings.Split(s, ",") {
		str = strings.Trim(str, " ")
		l = append(l, str)
	}

	return l
}

func string2bool(s string) bool {

	possible := []string{
		"true",
		"y",
		"yes",
	}

	b := false

	for _, p := range possible {

		if strings.ToLower(s) == p {
			b = true
			break
		}
	}

	return b
}

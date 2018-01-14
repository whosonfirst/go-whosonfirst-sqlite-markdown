package search

import (
	"github.com/whosonfirst/go-whosonfirst-markdown"
	"gopkg.in/russross/blackfriday.v2"
	"io"
	"log"
	"net/url"
	"strings"
)

type Indexer interface {
	IndexDocument(*markdown.Document) (*SearchDocument, error)
	Query(*SearchQuery) (interface{}, error) // not sure what this should really return, re-inflating a SearchDocument seems like overkill
	Close() error
}

type SearchDocument struct {
	Id       string
	Title    string
	Category string
	Tags     []string
	Authors  []string
	Date     string
	Links    map[string]*url.URL
	Images   map[string]int
	Body     []string
	Code     []string
}

type SearchQuery struct {
	QueryString string
}

func NewDefaultSearchQuery(qs string) (*SearchQuery, error) {

	q := SearchQuery{
		QueryString: qs,
	}

	return &q, nil
}

func NewSearchDocument(doc *markdown.Document) (*SearchDocument, error) {

	fm := doc.FrontMatter
	body := doc.Body

	links := make(map[string]*url.URL)
	images := make(map[string]int)

	search_doc := SearchDocument{
		Id:       fm.Permalink,
		Title:    fm.Title,
		Category: fm.Category,
		Tags:     fm.Tags,
		Authors:  fm.Authors,
		Date:     fm.Date,
		Body:     []string{},
		Code:     []string{},
		Images:   images,
		Links:    links,
	}

	params := blackfriday.HTMLRendererParameters{}
	renderer := blackfriday.NewHTMLRenderer(params)

	r := SearchRenderer{
		bf:  renderer,
		doc: &search_doc,
	}

	blackfriday.Run(body.Bytes(), blackfriday.WithRenderer(&r))

	return &search_doc, nil
}

type SearchRenderer struct {
	bf  *blackfriday.HTMLRenderer
	doc *SearchDocument
}

func (r *SearchRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {

	str_value := string(node.Literal)

	switch node.Type {
	case blackfriday.Text:

		str_value = strings.Trim(str_value, " ")

		if str_value != "" {

			if node.Parent.Type == blackfriday.Link {

				str_link := string(node.LinkData.Destination)
				link, err := url.Parse(str_link)

				if err == nil {
					r.doc.Links[str_link] = link
				}

			}

			// log.Println("TEXT", str_value)
			r.doc.Body = append(r.doc.Body, str_value)
		}
	case blackfriday.Softbreak:
		// pass
	case blackfriday.Hardbreak:
		// pass
	case blackfriday.Emph:
		// pass
	case blackfriday.Strong:
		// pass
	case blackfriday.Del:
		// pass
	case blackfriday.HTMLSpan:
		// pass
	case blackfriday.Link:

		if entering {
			str_link := string(node.LinkData.Destination)
			link, err := url.Parse(str_link)

			if err == nil {
				r.doc.Links[str_link] = link
			}
		}

	case blackfriday.Image:

		if entering {

			str_link := string(node.LinkData.Destination)
			link, err := url.Parse(str_link)

			if err == nil {
				r.doc.Links[str_link] = link
			}

		}
		// pass
	case blackfriday.Code:
		r.doc.Code = append(r.doc.Code, str_value)
	case blackfriday.Document:
		break
	case blackfriday.Paragraph:
		// pass
	case blackfriday.BlockQuote:
		// pass
	case blackfriday.HTMLBlock:
		// pass
	case blackfriday.Heading:
		// pass
	case blackfriday.HorizontalRule:
		// pass
	case blackfriday.List:
		// pass
	case blackfriday.Item:
		// pass
	case blackfriday.CodeBlock:
		r.doc.Code = append(r.doc.Code, str_value)
	case blackfriday.Table:
		// pass
	case blackfriday.TableCell:
		// pass
	case blackfriday.TableHead:
		// pass
	case blackfriday.TableBody:
		// pass
	case blackfriday.TableRow:
		// pass
	default:
		log.Println("Unknown node type " + node.Type.String())
	}
	return blackfriday.GoToNext
}

func (r *SearchRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	return
}

func (r *SearchRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	return
}

package markdown

import (
	"bytes"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-markdown/jekyll"
)

type Document struct {
	FrontMatter *jekyll.FrontMatter
	Body        *Body
}

func NewDocument(fm *jekyll.FrontMatter, body *Body) (*Document, error) {

	doc := Document{
		FrontMatter: fm,
		Body:        body,
	}

	return &doc, nil
}

type Body struct {
	*bytes.Buffer
}

func (b *Body) String() string {
	return fmt.Sprintf("%s", b.Bytes())
}

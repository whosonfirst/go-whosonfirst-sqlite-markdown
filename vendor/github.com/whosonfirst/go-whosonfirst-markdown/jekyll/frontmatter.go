package jekyll // https://jekyllrb.com/docs/frontmatter/

import (
	"bufio"
	"bytes"
	"log"
	"text/template"
)

var fm_template *template.Template

func init() {

	tm := `---
layout: {{ .Layout }}
permalink: {{ .Permalink }}
published: {{ .Published }}
title: {{ .Title }}
date: {{ .Date }}
category: {{ .Category}}
excerpt: {{ .Excerpt }}
authors: {{ .Authors }}
image: {{ .Image }}
tags: {{ .Tags }}
---`

	t, err := template.New("frontmatter").Parse(tm)

	if err != nil {
		log.Fatal(err)
	}

	fm_template = t
}

type FrontMatter struct {
	// default
	Layout    string
	Permalink string
	Published bool
	// out-of-the-box
	Category string
	Date     string
	// custom
	Title   string
	Excerpt string
	Image   string
	Authors []string
	Tags    []string
}

func (fm *FrontMatter) String() string {

	var b bytes.Buffer
	wr := bufio.NewWriter(&b)

	err := fm_template.Execute(wr, fm)

	if err != nil {
		return err.Error()
	}

	wr.Flush()
	return string(b.Bytes())
}

func EmptyFrontMatter() *FrontMatter {

	fm := FrontMatter{
		Title:     "",
		Excerpt:   "",
		Image:     "",
		Layout:    "",
		Category:  "",
		Published: false,
		Authors:   make([]string, 0),
		Tags:      make([]string, 0),
		Date:      "",
		Permalink: "",
	}

	return &fm
}

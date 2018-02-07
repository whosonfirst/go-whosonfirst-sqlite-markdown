package render

import (
	"text/template"
)

type FeedOptions struct {
	Format    string
	Input     string
	Output    string
	Items     int
	Templates *template.Template
}

func DefaultFeedOptions() *FeedOptions {

	opts := FeedOptions{
		Input:     "index.md",
		Format:    "rss_20",
		Output:    "rss_20.xml",
		Items:     10,
		Templates: nil,
	}

	return &opts
}

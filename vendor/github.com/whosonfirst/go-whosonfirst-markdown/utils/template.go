package utils

import (
	"html/template"
	"io/ioutil"
	"os"
)

func LoadTemplate(path string, name string) (*template.Template, error) {

	fh, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	// https://play.golang.org/p/V94BPN0uKD

	var fns = template.FuncMap{
		"plus1": func(x int) int {
			return x + 1
		},
	}

	return template.New(name).Funcs(fns).Parse(string(body))
}

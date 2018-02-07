package utils

import (
	"github.com/facebookgo/atomicfile"
	"github.com/whosonfirst/go-whosonfirst-markdown/render"
	"io"
	"log"
	"os"
	"path/filepath"
)

func WriteHTML(fh io.ReadCloser, root string, opts *render.HTMLOptions) error {

	if opts.Output == "STDOUT" {
		_, err := io.Copy(os.Stdout, fh)
		return err

	}

	index := filepath.Join(root, opts.Output)

	out, err := atomicfile.New(index, os.FileMode(0644))

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, fh)

	if err != nil {
		out.Abort()
		return err
	}

	log.Printf("wrote %s", index)
	return nil
}

func WriteFeed(fh io.ReadCloser, root string, opts *render.FeedOptions) error {

	if opts.Output == "STDOUT" {
		_, err := io.Copy(os.Stdout, fh)
		return err

	}

	index := filepath.Join(root, opts.Output)

	out, err := atomicfile.New(index, os.FileMode(0644))

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, fh)

	if err != nil {
		out.Abort()
		return err
	}

	log.Printf("wrote %s", index)
	return nil
}

package writer

import (
	"github.com/facebookgo/atomicfile"
	"io"
	_ "log"
	"os"
	"path/filepath"
)

type FSWriter struct {
	Writer
	root string
}

func NewFSWriter(root string) (Writer, error) {

	abs_root, err := filepath.Abs(root)

	if err != nil {
		return nil, err
	}

	w := FSWriter{
		root: abs_root,
	}

	return &w, nil
}

func (w *FSWriter) Write(rel_path string, fh io.ReadCloser) error {

	out_path := filepath.Join(w.root, rel_path)
	out_root := filepath.Dir(out_path)

	_, err := os.Stat(out_root)

	if os.IsNotExist(err) {

		err = os.MkdirAll(out_root, 0755)

		if err != nil {
			return err
		}
	}

	out, err := atomicfile.New(out_path, os.FileMode(0644))

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, fh)

	if err != nil {
		out.Abort()
		return err
	}

	return nil
}

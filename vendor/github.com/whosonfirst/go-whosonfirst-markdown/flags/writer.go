package flags

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-markdown/writer"
	"strings"
)

type WriterFlags []writer.Writer

func (fl *WriterFlags) String() string {
	return fmt.Sprintf("%v", *fl)
}

func (fl *WriterFlags) Set(value string) error {

	value = strings.Trim(value, " ")
	kv := strings.Split(value, "=")

	var wr writer.Writer
	var err error

	switch kv[0] {

	case "fs":

		if len(kv) < 2 {
			err = errors.New("Missing path")
		} else {
			path := kv[1]
			wr, err = writer.NewFSWriter(path)
		}
	case "null":
		wr, err = writer.NewNullWriter()
	case "stdout":
		wr, err = writer.NewStdoutWriter()
	default:
		err = errors.New("Invalid writer")
	}

	if err != nil {
		return err
	}

	*fl = append(*fl, wr)
	return nil
}

func (fl *WriterFlags) ToWriter() (writer.Writer, error) {
	return writer.NewMultiWriter(*fl...)
}

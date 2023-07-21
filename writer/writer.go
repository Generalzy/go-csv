package writer

import (
	"encoding/csv"
	gocsv "github.com/generalzy/go-csv"
	"os"
)

type Writer struct {
	wt         *csv.Writer
	fp         *os.File
	headLength int
	head       []string
}

func NewWriter(filename string) (*Writer, error) {
	fp, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &Writer{wt: csv.NewWriter(fp), fp: fp}, nil
}

func (w *Writer) Close() error {
	w.wt.Flush()

	w.head = nil
	w.wt = nil

	return w.fp.Close()
}

func (w *Writer) WriteHead(head []string) error {
	w.head = head
	w.headLength = len(head)
	return w.wt.Write(head)
}

func (w *Writer) WriteLine(line []string) error {
	if w.headLength == 0 {
		return gocsv.MissingHeadError
	}
	return w.wt.Write(line)
}

func (w *Writer) WriteLines(lines [][]string) error {
	if w.headLength == 0 {
		return gocsv.MissingHeadError
	}

	return w.wt.WriteAll(lines)
}

func (w *Writer) Head() []string {
	return w.head
}

func (w *Writer) Scope() int {
	return w.headLength
}

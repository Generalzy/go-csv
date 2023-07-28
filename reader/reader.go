package reader

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
)

type Reader struct {
	rd         *csv.Reader
	fp         *os.File
	headLength int
	head       []string
	Info       os.FileInfo
}

func NewReader(filename string) (*Reader, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	stat, err := fp.Stat()
	if err != nil {
		return nil, err
	}

	rd := csv.NewReader(fp)
	rd.TrimLeadingSpace = true
	
	return &Reader{rd: rd, head: make([]string, 0, 0), fp: fp, Info: stat}, nil
}

func (r *Reader) Close() error {
	r.rd = nil
	r.head = nil
	return r.fp.Close()
}

func (r *Reader) ReadHead() ([]string, error) {
	if r.headLength == 0 {
		head, err := r.rd.Read()
		if err != nil {
			return nil, err
		}

		r.head = head
		r.headLength = len(head)
		return head, nil
	}

	return r.head, nil
}

func (r *Reader) ReadLine() ([]string, error) {
	return r.rd.Read()
}

func (r *Reader) ReadLines() ([][]string, error) {
	return r.rd.ReadAll()
}

func (r *Reader) ReadWithFn(fn func(line []string) error) error {
	for {
		line, err := r.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		if err = fn(line); err != nil {
			return err
		}
	}
	return nil
}

func (r *Reader) Head() []string {
	return r.head
}

func (r *Reader) Scope() int {
	return r.headLength
}

func (r *Reader) SetComma(comma rune) {
	r.rd.Comma = comma
}

func (r *Reader) SetComment(comment rune) {
	r.rd.Comment = comment
}

func (r *Reader) SetLazyQuotes(fl bool) {
	// default is false
	r.rd.LazyQuotes = fl
}

func (r *Reader) SetTrimLeadingSpace(fl bool) {
	// default is true
	r.rd.TrimLeadingSpace = fl
}

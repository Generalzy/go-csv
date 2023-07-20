package reader

import (
	"encoding/csv"
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

	return &Reader{rd: csv.NewReader(fp), head: make([]string, 0, 0), fp: fp, Info: stat}, nil
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

func (r *Reader) Head() []string {
	return r.head
}

func (r *Reader) Scope() int {
	return r.headLength
}

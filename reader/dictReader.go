package reader

import (
	"errors"
	"io"
	"os"
)

type DictReader struct {
	rd   *Reader
	Info os.FileInfo
}

func NewDictReader(filename string) (*DictReader, error) {
	rd, err := NewReader(filename)
	if err != nil {
		return nil, err
	}

	if _, err = rd.ReadHead(); err != nil {
		return nil, err
	}

	return &DictReader{rd: rd, Info: rd.Info}, nil
}

func (d *DictReader) readHead() ([]string, error) {
	return d.rd.ReadHead()
}

func (d *DictReader) Close() error {
	return d.rd.Close()
}

func (d *DictReader) ReadLine() (map[string]string, error) {
	dictLine := make(map[string]string, d.rd.headLength)

	line, err := d.rd.ReadLine()
	if err != nil {
		return nil, err
	}

	for filedIndex, fieldVal := range line {
		dictLine[d.rd.head[filedIndex]] = fieldVal
	}

	return dictLine, nil
}

func (d *DictReader) ReadLines() ([]map[string]string, error) {
	dictLines := make([]map[string]string, 0, 0)

	for {
		dictLine, err := d.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		dictLines = append(dictLines, dictLine)
	}

	return dictLines, nil
}

func (d *DictReader) Head() []string {
	return d.rd.Head()
}

func (d *DictReader) Scope() int {
	return d.rd.Scope()
}

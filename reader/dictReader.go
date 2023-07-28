package reader

import (
	"encoding/json"
	"errors"
	gocsv "github.com/generalzy/go-csv"
	"io"
	"os"
	"reflect"
	"strconv"
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

func (d *DictReader) Close() error {
	return d.rd.Close()
}

func (d *DictReader) readHead() ([]string, error) {
	return d.rd.ReadHead()
}

func (d *DictReader) ReadDictLine() (map[string]string, error) {
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

func (d *DictReader) ReadDictLines() ([]map[string]string, error) {
	dictLines := make([]map[string]string, 0, 0)

	for {
		dictLine, err := d.ReadDictLine()
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

func (d *DictReader) ReadDictWith(fn func(dictLine map[string]string) error) error {
	for {
		dictLine, err := d.ReadDictLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		if err = fn(dictLine); err != nil {
			return err
		}
	}
	return nil
}

func (d *DictReader) Head() []string {
	return d.rd.Head()
}

func (d *DictReader) Scope() int {
	return d.rd.Scope()
}

func (d *DictReader) BindWithStruct(dst interface{}) error {
	v := reflect.ValueOf(dst)

	// Check if dst is a pointer to a struct.
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return gocsv.InvalidTypeError
	}

	// Get the struct type from the interface.
	t := v.Elem().Type()

	// Read line
	m, err := d.ReadDictLine()
	if err != nil {
		return err
	}

	// Loop through the struct fields.
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Get the "csv" tag.
		tag := field.Tag.Get(gocsv.Tag)

		// Skip fields without the "csv" tag.
		if tag == "" {
			continue
		}

		var value reflect.Value

		switch field.Type.Kind() {
		case reflect.String:
			value = reflect.ValueOf(m[tag])
		case reflect.Int:
			intValue, _ := strconv.Atoi(m[tag])
			value = reflect.ValueOf(intValue)
		case reflect.Float64 | reflect.Float32:
			floatValue, _ := strconv.ParseFloat(m[tag], 64)
			value = reflect.ValueOf(floatValue)
		default:
			return gocsv.UnsupportedTypeError
		}
		// Extract the value using reflection.

		// Set the value in the struct field.
		dstField := v.Elem().Field(i)

		if !dstField.CanSet() {
			return gocsv.CannotSetError
		}

		dstField.Set(value)
	}

	return nil
}

func (d *DictReader) BindWithJson() ([]byte, error) {
	line, err := d.ReadDictLine()
	if err != nil {
		return nil, err
	}

	return json.Marshal(line)
}

func (d *DictReader) SetComma(comma rune) {
	d.rd.SetComma(comma)
}

func (d *DictReader) SetComment(comment rune) {
	d.rd.SetComment(comment)
}

func (d *DictReader) SetLazyQuotes(fl bool) {
	// default is false
	d.rd.SetLazyQuotes(fl)
}

func (d *DictReader) SetTrimLeadingSpace(fl bool) {
	// default is true
	d.rd.SetTrimLeadingSpace(fl)
}

package writer

import (
	"errors"
	"fmt"
	gocsv "github.com/genralzy/go-csv"
	"reflect"
	"strconv"
)

type DictWriter struct {
	wt *Writer
}

func NewDictWriter(filename string) (*DictWriter, error) {
	wt, err := NewWriter(filename)
	if err != nil {
		return nil, err
	}
	return &DictWriter{wt: wt}, nil
}

func (d *DictWriter) WriteHead(head []string) error {
	return d.wt.WriteHead(head)
}

func (d *DictWriter) Close() error {
	return d.wt.Close()
}

func (d *DictWriter) writeDictLine(line map[string]string) error {
	lineSlice := make([]string, d.wt.headLength)

	for index, field := range d.wt.head {
		lineSlice[index] = line[field]
	}

	return d.wt.WriteLine(lineSlice)
}

func (d *DictWriter) WriteDictLines(lines []map[string]string) error {
	for _, line := range lines {
		if err := d.writeDictLine(line); err != nil {
			return err
		}
	}

	return nil
}

func (d *DictWriter) WriteLine(line interface{}) error {
	v := reflect.ValueOf(line)

	switch v.Kind() {
	case reflect.Map:
		dictLine, _ := line.(map[string]string)
		return d.writeDictLine(dictLine)
	case reflect.Slice:
		sliceLine, _ := line.([]string)
		return d.wt.WriteLine(sliceLine)
	case reflect.Pointer:
		if v.IsNil() {
			return errors.New("nil pointer")
		}
		return d.WriteLine(v.Elem().Interface())
	case reflect.Struct:
		t := v.Type()

		l := make(map[string]string, d.wt.headLength)

		for fieldIndex := 0; fieldIndex < v.NumField(); fieldIndex++ {
			field := t.Field(fieldIndex)

			fieldTagName := field.Tag.Get(gocsv.Tag)
			fieldVal := v.FieldByName(field.Name)

			var fieldValStr string

			switch fieldVal.Kind() {
			case reflect.String:
				fieldValStr = fieldVal.String()
			case reflect.Int:
				fieldValStr = strconv.FormatInt(fieldVal.Int(), 10)
			default:
				// Handle other types accordingly.
				return fmt.Errorf("unsupported type for field %s", field.Name)
			}
			l[fieldTagName] = fieldValStr
		}

		return d.writeDictLine(l)
	default:
		return errors.New("error type")
	}
}

func (d *DictWriter) Scope() int {
	return d.wt.Scope()
}

func (d *DictWriter) Head() []string {
	return d.wt.Head()
}

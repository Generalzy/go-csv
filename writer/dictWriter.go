package writer

import (
	gocsv "github.com/generalzy/go-csv"
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

func (d *DictWriter) WriteDictLine(line map[string]string) error {
	lineSlice := make([]string, d.wt.headLength)

	for index, field := range d.wt.head {
		lineSlice[index] = line[field]
	}

	return d.wt.WriteLine(lineSlice)
}

func (d *DictWriter) WriteDictLines(lines []map[string]string) error {
	for _, line := range lines {
		if err := d.WriteDictLine(line); err != nil {
			return err
		}
	}

	return nil
}

func (d *DictWriter) WriteLines(lines []interface{}) error {
	for _, line := range lines {
		if err := d.WriteLine(line); err != nil {
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
		return d.WriteDictLine(dictLine)
	case reflect.Slice:
		sliceLine, _ := line.([]string)
		return d.wt.WriteLine(sliceLine)
	case reflect.Pointer:
		if v.IsNil() {
			return gocsv.NilPointerError
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
			case reflect.Float64 | reflect.Float32:
				fieldValStr = strconv.FormatFloat(fieldVal.Float(), 'f', 2, 64)
			default:
				// Handle other types accordingly.
				return gocsv.UnsupportedTypeError
			}
			l[fieldTagName] = fieldValStr
		}

		return d.WriteDictLine(l)
	default:
		return gocsv.InvalidTypeError
	}
}

func (d *DictWriter) Scope() int {
	return d.wt.Scope()
}

func (d *DictWriter) Head() []string {
	return d.wt.Head()
}

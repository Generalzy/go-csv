package utils

import (
	"bufio"
	"encoding/json"
	"github.com/generalzy/go-csv/reader"
	"io"
	"os"
)

func Copy(src, dst string) error {
	srcFp, err := os.Open(src)
	if err != nil {
		return err
	}
	dstFp, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(dstFp, srcFp)
	return err
}

func Append(src, dst string) error {
	srcFp, err := os.Open(src)
	if err != nil {
		return err
	}

	srcReader := bufio.NewReader(srcFp)
	// read head
	_, err = srcReader.ReadSlice('\n')

	if err != nil {
		return err
	}

	dstFp, err := os.OpenFile(dst, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	_, err = io.Copy(dstFp, srcReader)
	return err
}

func CsvToJson(src, dst string) error {
	dstFp, err := os.Create(dst)
	if err != nil {
		return err
	}

	jsonEncoder := json.NewEncoder(dstFp)

	rd, err := reader.NewDictReader(src)
	if err != nil {
		return err
	}

	if err = rd.ReadDictWith(func(dictLine map[string]string) error {
		return jsonEncoder.Encode(dictLine)
	}); err != nil {
		return err
	}

	return nil
}

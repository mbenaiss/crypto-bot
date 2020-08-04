package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
)

type Client struct {
	delimiter rune
	columns   []string
}

func New(delimiter rune, columns []string) *Client {
	return &Client{
		delimiter: delimiter,
		columns:   columns,
	}
}

func (c *Client) Read(filename string) ([]map[string]string, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %w", filename, err)
	}

	r := csv.NewReader(bytes.NewReader(f))
	r.Comma = c.delimiter

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	columns := make(map[string]int)
	for i, col := range records[0] {
		columns[col] = i
	}

	rColumns := make(map[int]string, len(c.columns))
	for n, i := range columns {
		if len(c.columns) == 0 {
			rColumns[i] = n
			continue
		}
		for _, v := range c.columns {
			if n == v {
				rColumns[i] = n
			}
		}
	}

	result := make([]map[string]string, 0, len(records))

	for i, rows := range records {
		if i > 1 {
			r := map[string]string{}
			for j, row := range rows {
				col, ok := rColumns[j]
				if ok {
					r[col] = row
					result = append(result, r)
				}
			}
		}
	}

	return result, nil
}

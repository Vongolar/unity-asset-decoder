package decode

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

func GetGUID(r io.Reader) (string, error) {
	br := bufio.NewReader(r)

	var line []byte
	var err error

	prex := "guid:"
	lenPrex := len(prex)

	for err == nil {
		line, _, err = br.ReadLine()

		if err != nil {
			break
		}

		lineStr := string(bytes.TrimSpace(line))
		if lineStr[:lenPrex] == prex {
			return strings.TrimSpace(lineStr[lenPrex:]), nil
		}
	}

	return "", err
}

package yaml

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

const yamlSeparator = "\n---"

// splitYAMLDocument is a bufio.SplitFunc for splitting YAML streams into individual documents.
func splitYAMLDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	sep := len([]byte(yamlSeparator))
	if i := bytes.Index(data, []byte(yamlSeparator)); i >= 0 {
		// We have a potential document terminator
		i += sep
		after := data[i:]
		if len(after) == 0 {
			// we can't read any more characters
			if atEOF {
				return len(data), data[:len(data)-sep], nil
			}
			return 0, nil, nil
		}
		if j := bytes.IndexByte(after, '\n'); j >= 0 {
			return i + j + 1, data[0 : i-sep], nil
		}
		return 0, nil, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func Split(cf string) ([]*bytes.Buffer, error) {
	fd, err := os.Open(cf)
	if err != nil {
		return nil, err
	}
	return SplitReader(fd)
}

func SplitReader(rd io.Reader) ([]*bytes.Buffer, error) {
	bs := []*bytes.Buffer{}
	scanner := bufio.NewScanner(rd)
	scanner.Split(splitYAMLDocument)
	for scanner.Scan() {
		bs = append(bs, bytes.NewBuffer(scanner.Bytes()))
	}
	return bs, nil
}

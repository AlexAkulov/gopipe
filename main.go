package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type document struct {
	Fields map[string]interface{} `json:"fields"`
	Source map[string]interface{} `json:"_source"`
}

func convert(sourceLine []byte) ([]byte, error) {
	var doc document
	if err := json.Unmarshal(sourceLine, &doc); err != nil {
		return nil, err
	}
	if _, ok := doc.Fields["_timestamp"]; !ok {
		return nil, fmt.Errorf("BAD DOC [%s]", sourceLine)

	}
	doc.Source["@timestamp"] = doc.Fields["_timestamp"]
	delete(doc.Fields, "_timestamp")
	return json.Marshal(doc)
}

func main() {
	_, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(os.Stdin)
	var (
		line                 bytes.Buffer
		badLines, totalLines int
	)
	for {
		part, isPrefix, err := reader.ReadLine()
		if err == nil {
			line.Write(part)
			if !isPrefix {
				totalLines++
				newLine, err := convert(line.Bytes())
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					badLines++
				} else {
					l := make([]byte, line.Len())
					copy(l, line.Bytes())
					fmt.Println(string(newLine))
				}
				line.Reset()
			}
			continue
		}
		if err != io.EOF {
			panic(err)
		}
		break
	}
}

package internal

import (
	"bufio"
	"encoding/xml"
	"io"
	"log"
	"os"
	"strings"
)

type DefaultFileReader struct{}

func (d DefaultFileReader) Read(path string) (<-chan string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	out := make(chan string)
	go func() {
		defer file.Close()
		defer close(out)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			out <- scanner.Text()
		}
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
	}()
	return out, nil

}

type XMLFileReader struct{}

func (x XMLFileReader) Read(path string) (<-chan string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	out := make(chan string)
	go func() {
		defer file.Close()
		defer close(out)
		decoder := xml.NewDecoder(file)

		var inBody bool
		for {
			token, err := decoder.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				return
			}
			switch t := token.(type) {
			case xml.StartElement:
				if t.Name.Local == "body" {
					inBody = true
				}
			case xml.EndElement:
				if t.Name.Local == "body" {
					inBody = false
				}
			case xml.CharData:
				if inBody {
					out <- strings.ToLower(strings.TrimSpace(string(t)))
				}
			}
		}
	}()
	return out, nil

}

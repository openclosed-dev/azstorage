package blob

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type directoryHandler interface {
	handleDirectory(path string)
}

type listParser struct {
	file    *os.File
	scanner *bufio.Scanner
	handler directoryHandler
}

func newListParser(path string, handler directoryHandler) (*listParser, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open list file \"%s\": %v", path, err)
		return nil, err
	}

	return &listParser{
		file:    file,
		scanner: bufio.NewScanner(file),
		handler: handler,
	}, nil
}

func (p *listParser) close() {
	p.file.Close()
}

func (p *listParser) parseAll() error {
	var scanner = p.scanner
	var handler = p.handler
	var line string
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		// Skips comment
		if strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "/")
		if !strings.HasSuffix(line, "/") {
			line = line + "/"
		}
		if len(line) > 0 {
			handler.handleDirectory(line)
		}
	}

	var err = scanner.Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error has occurred while parsing the list file: %v", err)
		return err
	}

	return nil
}

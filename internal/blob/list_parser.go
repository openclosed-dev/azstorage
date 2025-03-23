package blob

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const logPrefix = "[parser]"

type listEntryHandler interface {
	handleBlob(path string)
	handleDirectory(path string)
}

type listParser struct {
	file    *os.File
	scanner *bufio.Scanner
	handler listEntryHandler
}

func newListParser(path string, handler listEntryHandler) (*listParser, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%s failed to open the list file \"%s\": %v", logPrefix, path, err)
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

	log.Printf("%s Start parsing list file \"%s\"", logPrefix, p.file.Name())

	var scanner = p.scanner
	var handler = p.handler
	var line string
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		// Skips comment
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		if len(line) > 1 {
			line = strings.TrimPrefix(line, "/")
		}

		if strings.HasSuffix(line, "/") {
			handler.handleDirectory(line)
		} else {
			handler.handleBlob(line)
		}
	}

	var err = scanner.Err()
	if err != nil {
		return fmt.Errorf("error has occurred while parsing the list file: %w", err)
	}

	log.Printf("%s Finished parsing list file \"%s\"", logPrefix, p.file.Name())

	return nil
}

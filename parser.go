package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sync"
)

func readAllComponents(in io.ReadCloser, out chan component, wait *sync.WaitGroup) {
	scanner := bufio.NewScanner(bufio.NewReader(in))
	var lines []string
	var componenetFound bool
	var inComponent bool
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == 'C' {
			lines = make([]string, 0, 10)
			componenetFound = true
			lines = append(lines, line)
		} else if componenetFound {
			if line[0] == '{' {
				inComponent = true
			} else if !inComponent {
				log.Printf("Found component without attributes. Skipping. (%s)", lines[0])
				inComponent = false
				componenetFound = false
			} else if line[0] == '}' {
				inComponent = false
				componenetFound = false
				comp := NewComponent(lines)
				out <- comp
			} else {
				lines = append(lines, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	in.Close()
	wait.Done()
}

func startParserOnIOReaders(in []io.ReadCloser) chan component {
	wait := new(sync.WaitGroup)
	out := make(chan component)
	for _, i := range in {
		wait.Add(1)
		go readAllComponents(i, out, wait)
	}

	go func() {
		wait.Wait()
		close(out)
	}()

	return out
}

func ParseFiles(inFileNames []string) componentArray {
	inReaders := make([]io.ReadCloser, 0, len(inFileNames))
	if len(inFileNames) == 0 {
		inReaders = append(inReaders, os.Stdin)
	} else {
		for _, fileName := range inFileNames {
			inFile, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			} else {
				inReaders = append(inReaders, inFile)
			}
		}
	}

	outChan := startParserOnIOReaders(inReaders)
	components := NewComponentArray(100)
	for c := range outChan {
		components = append(components, c)
	}

	return components
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"push-swap-go/internal/pushswap"
)

type filePair struct {
	Input  string
	Output string
}

type filePairs []filePair

// String is required by the flag.Value interface.
func (f *filePairs) String() string {
	var pairs []string

	for _, pair := range *f {
		pairs = append(pairs, fmt.Sprintf("%s,%s", pair.Input, pair.Output))
	}

	return strings.Join(pairs, " ")
}

// custom parsing logic for `filePairs`.
func (f *filePairs) Set(value string) error {
	parts := strings.Split(value, ",")
	pair := filePair{Input: parts[0]}

	pair.Output = pair.Input + ".output"
	if len(parts) > 1 {
		pair.Output = parts[1]
	}

	*f = append(*f, pair)
	return nil
}

func printHelp() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "\tReads space separated numbers from stdin and writes the push-swap instructions")
	fmt.Fprintf(flag.CommandLine.Output(), "\tfor sorting them to stdout by default.")
	fmt.Fprintln(flag.CommandLine.Output())
	flag.PrintDefaults()
}

func readNumbers(file string, allowDups bool) ([]float64, error) {
	input := os.Stdin

	if file != "-" {
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("opening file: %v", err)
		}

		defer f.Close()
		input = f
	}

	var numStrings []string
	var line strings.Builder
	inputReader := bufio.NewReader(input)

	for {
		buf, isPrefix, err := inputReader.ReadLine()

		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("reading file: %v", err)
		}

		// Process the line if it has content
		if len(buf) > 0 {
			line.WriteString(string(buf))

			if !isPrefix {
				numStrings = append(numStrings, strings.Fields(line.String())...)
				line.Reset()
			}
		}

		// Exit the loop on EOF
		if err == io.EOF {
			break
		}
	}

	numbers, err := pushswap.ParseNumberSlice(numStrings, allowDups)
	if err != nil {
		return nil, err
	}

	return numbers, nil
}

func writeInstructions(file string, instructions []pushswap.Operation) (int, error) {
	output := os.Stdout

	if file != "-" {
		f, err := os.Create(file)
		if err != nil {
			return -1, fmt.Errorf("opening file: %v", err)
		}

		defer f.Close()
		output = f
	}

	bytesWritten := 0
	for _, op := range instructions {
		n, err := fmt.Fprintln(output, op)
		if err != nil {
			return n, fmt.Errorf("writing to file: %v", err)
		}

		bytesWritten += n
	}

	return bytesWritten, nil
}

func main() {
	allowDups := flag.Bool("allow-duplicates", false, "allow duplicate values in the number list")
	var files filePairs

	flag.Var(&files, "files", "specifies an input file and an optional output file separated by a comma.")
	flag.Usage = printHelp
	flag.Parse()

	if len(files) < 1 {
		files = append(files, filePair{Input: "-", Output: "-"})
	}

	for _, pair := range files {

		numbers, err := readNumbers(pair.Input, *allowDups)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		instructions := pushswap.TurkAlgorithm(numbers)

		_, err = writeInstructions(pair.Output, instructions)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}
	}
}

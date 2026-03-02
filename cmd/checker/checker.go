package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"push-swap-go/internal/pushswap"
)

type filePair struct {
	instructionsFile string
	numbersFile      string
}

type filePairs []filePair

// String is required by the flag.Value interface.
func (f *filePairs) String() string {
	var pairs []string

	for _, pair := range *f {
		pairs = append(pairs, fmt.Sprintf("%s,%s", pair.instructionsFile, pair.numbersFile))
	}

	return strings.Join(pairs, " ")
}

// custom parsing logic for `filePairs`.
func (f *filePairs) Set(value string) error {
	parts := strings.Split(value, ",")

	if len(parts) != 2 {
		return fmt.Errorf("usage: instructions_file,numbers_file")
	}

	pair := filePair{instructionsFile: parts[0], numbersFile: parts[1]}
	*f = append(*f, pair)
	return nil
}

func printHelp() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] numbers...", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "\tReads a list of instructions from stdin to sort a list of numbers given")
	fmt.Fprintf(flag.CommandLine.Output(), "\tvia the command line by default.")
	fmt.Fprintln(flag.CommandLine.Output())
	flag.PrintDefaults()
}

func readNumbers(file string, allowDups bool) ([]float64, error) {
	input, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("opening file: %v", err)
	}
	defer input.Close()

	var numStrings []string
	inputScanner := bufio.NewScanner(input)

	for inputScanner.Scan() {
		numStrings = append(numStrings, strings.Fields(inputScanner.Text())...)
	}

	err = inputScanner.Err()
	if err != nil {
		return nil, fmt.Errorf("reading file: %v", err)
	}

	numbers, err := pushswap.ParseNumberSlice(numStrings, allowDups)
	if err != nil {
		return nil, err
	}

	return numbers, nil
}

func readInstructions(file string) ([]pushswap.Operation, error) {
	var input *os.File

	if file == "-" {
		input = os.Stdin
	} else {
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("opening file: %v", err)
		}
		defer f.Close()
		input = f
	}

	var instructions []pushswap.Operation
	inputScanner := bufio.NewScanner(input)

	for inputScanner.Scan() {
		op := pushswap.Operation(strings.Fields(inputScanner.Text())[0])

		switch op {
		case pushswap.Invalid:
		case pushswap.PA, pushswap.PB, pushswap.RA,
			pushswap.RB, pushswap.RR, pushswap.RRA,
			pushswap.RRB, pushswap.RRR, pushswap.SA,
			pushswap.SB, pushswap.SS:
			instructions = append(instructions, pushswap.Operation(op))
		default:
			return nil, fmt.Errorf("unrecognised command %q", op)
		}
	}

	err := inputScanner.Err()
	if err != nil {
		return nil, fmt.Errorf("reading file: %v", err)
	}

	return instructions, nil
}

func checkStacks(ds pushswap.DoubleStack[float64], reference []float64) (string, error) {
	if ds.B.Len() > 0 || ds.A.Len() != len(reference) {
		return "", fmt.Errorf("Got:\n%v\nExpected:\n%v", ds, reference)
	}

	for i, val := range ds.A.All() {
		if val != reference[i] {
			return "KO", fmt.Errorf("Got:\n%v\nExpected:\n%v", ds, reference)
		}
	}

	return "OK", nil
}

func main() {
	allowDups := flag.Bool("allow-duplicates", false, "allow duplicate values in the number list")
	var files filePairs

	flag.Var(&files, "files", "specifies an instructions file and a numbers file separated by a comma.")
	flag.Usage = printHelp
	flag.Parse()
	args := flag.Args()

	if len(files) < 1 {
		if len(args) < 1 {
			printHelp()
			os.Exit(1)
		}

		var numStrings []string
		for _, a := range args {
			numStrings = append(numStrings, strings.Fields(a)...)
		}

		numbers, err := pushswap.ParseNumberSlice(numStrings, *allowDups)
		if err != nil {
			log.Fatalln("ERROR:", err)
		}

		ds := pushswap.NewDoubleStack(numbers...)
		instructions, err := readInstructions("-")
		if err != nil {
			log.Fatalln("ERROR:", err)
		}

		ds.ExecuteInstructions(instructions)
		sorted := make([]float64, len(numbers))

		copy(sorted, numbers)
		slices.Sort(sorted)
		status, err := checkStacks(*ds, sorted)
		if err != nil {
			log.Fatalln("ERROR:", err)
		}

		fmt.Println(status)
	} else {
		for _, pair := range files {
			numbers, err := readNumbers(pair.numbersFile, *allowDups)
			if err != nil {
				log.Println("ERROR:", err)
				continue
			}

			ds := pushswap.NewDoubleStack(numbers...)
			instructions, err := readInstructions(pair.instructionsFile)
			if err != nil {
				log.Println("ERROR:", err)
				continue
			}

			ds.ExecuteInstructions(instructions)
			sorted := make([]float64, len(numbers))

			copy(sorted, numbers)
			slices.Sort(sorted)
			status, err := checkStacks(*ds, sorted)
			if err != nil {
				log.Println("ERROR:", err)
				continue
			}

			fmt.Println(status)
		}
	}
}

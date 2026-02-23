package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"push-swap-go/internal/pushswap"
)

func printHelp() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <number list>", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "\t<number list> is a single string of space separated numbers")
	fmt.Fprintln(flag.CommandLine.Output())
	flag.PrintDefaults()
}

func main() {
	allowDups := flag.Bool("allow-duplicates", false, "allow duplicate values in the number list")
	flag.Usage = printHelp

	flag.Parse()
	args := flag.Args()
	var numStrings []string

	for _, s := range args {
		numStrings = append(numStrings, strings.Fields(s)...)
	}

	nums, err := pushswap.ParseNumberSlice(numStrings, *allowDups)
	if err != nil {
		fmt.Fprintln(os.Stdout, "ERROR:", err)
		os.Exit(1)
	}

	instructions := pushswap.TurkAlgorithm(nums)

	for _, op := range instructions {
		fmt.Println(op)
	}
}

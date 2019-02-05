package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

var arguments = struct {
	Output      string
	StartID     int
	StopID      int
	Concurrency int
	Verbose     bool
}{}

func parseArgs(args []string) {
	// Create new parser object
	parser := argparse.NewParser("TenorArchiver", "Tenor.com archiver")

	// Create flags
	output := parser.String("o", "output", &argparse.Options{
		Required: false,
		Help:     "Output directory",
		Default:  "Downloads/"})

	startID := parser.Int("", "start-id", &argparse.Options{
		Required: false,
		Help:     "First ID to scrape",
		Default:  1})

	stopID := parser.Int("", "stop-id", &argparse.Options{
		Required: false,
		Help:     "Last ID to scrape",
		Default:  10000000000})

	concurrency := parser.Int("j", "concurrency", &argparse.Options{
		Required: false,
		Help:     "Concurrency",
		Default:  4})

	verbose := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Help:     "Verbose output",
		Default:  false})

	// Parse input
	err := parser.Parse(args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	// Fill arguments structure
	arguments.Output = *output
	arguments.Concurrency = *concurrency
	arguments.StartID = *startID
	arguments.StopID = *stopID
	arguments.Verbose = *verbose
}

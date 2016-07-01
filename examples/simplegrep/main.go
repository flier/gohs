/*
 * Hyperscan example program 1: simplegrep
 *
 * This is a simple example of Hyperscan's most basic functionality: it will
 * search a given input file for a pattern supplied as a command-line argument.
 * It is intended to demonstrate correct usage of the hs_compile and hs_scan
 * functions of Hyperscan.
 *
 * Patterns are scanned in 'DOTALL' mode, which is equivalent to PCRE's '/s'
 * modifier. This behaviour can be changed by modifying the "flags" argument to
 * hs_compile.
 *
 * Build instructions:
 *
 *     go build github.com/flier/gohs/examples/simplegrep
 *
 * Usage:
 *
 *     ./simplegrep <pattern> <input file>
 *
 * Example:
 *
 *     ./simplegrep int simplegrep.c
 *
 */
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/flier/gohs/hyperscan"
)

var (
	flagNoColor    = flag.Bool("C", false, "Disable colorized output.")
	flagByteOffset = flag.Bool("b", false, "The offset in bytes of a matched pattern is displayed in front of the respective matched line")
)

var theme = func(s string) string { return s }

func highlight(s string) string {
	return "\033[35m" + s + "\033[0m"
}

/**
 * This is the function that will be called for each match that occurs. @a ctx
 * is to allow you to have some application-specific state that you will get
 * access to for each match. In our simple example we're just going to use it
 * to pass in the pattern that was being searched for so we can print it out.
 */
func eventHandler(id uint, from, to uint64, flags uint, context interface{}) error {
	inputData := context.([]byte)

	start := bytes.LastIndexByte(inputData[:from], '\n')
	end := int(to) + bytes.IndexByte(inputData[to:], '\n')

	if start == -1 {
		start = 0
	} else {
		start += 1
	}

	if end == -1 {
		end = len(inputData)
	}

	if *flagByteOffset {
		fmt.Printf("%d: ", start)
	}

	fmt.Printf("%s%s%s\n", inputData[start:from], theme(string(inputData[from:to])), inputData[to:end])

	return nil
}

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <pattern> <input file>\n", os.Args[0])
		os.Exit(-1)
	}

	if !*flagNoColor {
		stat, _ := os.Stdout.Stat()

		if stat != nil && stat.Mode()&os.ModeType != 0 {
			theme = highlight
		}
	}

	pattern := hyperscan.NewPattern(flag.Arg(0), hyperscan.DotAll|hyperscan.SomLeftMost)
	inputFN := flag.Arg(1)

	/* First, we attempt to compile the pattern provided on the command line.
	 * We assume 'DOTALL' semantics, meaning that the '.' meta-character will
	 * match newline characters. The compiler will analyse the given pattern and
	 * either return a compiled Hyperscan database, or an error message
	 * explaining why the pattern didn't compile.
	 */
	database, err := hyperscan.NewBlockDatabase(pattern)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to compile pattern \"%s\": %s\n", pattern.String(), err.Error())
		os.Exit(-1)
	}

	defer database.Close()

	/* Next, we read the input data file into a buffer. */
	inputData, err := ioutil.ReadFile(inputFN)

	if err != nil {
		os.Exit(-1)
	}

	/* Finally, we issue a call to hs_scan, which will search the input buffer
	 * for the pattern represented in the bytecode. Note that in order to do
	 * this, scratch space needs to be allocated with the hs_alloc_scratch
	 * function. In typical usage, you would reuse this scratch space for many
	 * calls to hs_scan, but as we're only doing one, we'll be allocating it
	 * and deallocating it as soon as our matching is done.
	 *
	 * When matches occur, the specified callback function (eventHandler in
	 * this file) will be called. Note that although it is reminiscent of
	 * asynchronous APIs, Hyperscan operates synchronously: all matches will be
	 * found, and all callbacks issued, *before* hs_scan returns.
	 *
	 * In this example, we provide the input pattern as the context pointer so
	 * that the callback is able to print out the pattern that matched on each
	 * match event.
	 */
	scratch, err := hyperscan.NewScratch(database)

	if err != nil {
		fmt.Fprint(os.Stderr, "ERROR: Unable to allocate scratch space. Exiting.\n")
		os.Exit(-1)
	}

	defer scratch.Free()

	fmt.Printf("Scanning %d bytes with Hyperscan\n", len(inputData))

	if err := database.Scan(inputData, scratch, eventHandler, inputData); err != nil {
		fmt.Fprint(os.Stderr, "ERROR: Unable to scan input buffer. Exiting.\n")
		os.Exit(-1)
	}

	/* Scanning is complete, any matches have been handled, so now we just
	 * clean up and exit.
	 */

	return
}

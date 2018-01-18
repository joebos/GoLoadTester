package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

// how to run: ./goloadtest lcoalhost:10497.
// test result will be written into perf_test.txt file, and console.
func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}

	// load testing settings
	maxNumOfClient := 601
	step := 10

	// add test file and stdout as log output
	testResultFile := "perf_test.txt"
	logFile, _ := os.Create(testResultFile)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate)

	serverAddress := os.Args[1]

	loadtester := Loadtester{address: serverAddress}

	// start load testing by incrementing # of clients. Each client will continuously make 10 GET requests.
	// Line numbers are randomly generated.
	for i := 1; i <= maxNumOfClient; i = i + step {
		loadtester.StartLoadTest(i)
	}

}

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jasonlvhit/gocron"
	"github.com/jojomi/check"
)

func task() {
	fmt.Println("I am runnning task.")
}

func taskWithParams(a int, b string) {
	fmt.Println(a, b)
}

// CheckList is a list of checks
type CheckList []check.Check

var checkList CheckList
var cp = check.ConfigParser{}

func visit(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		return nil
	}
	if !strings.HasSuffix(f.Name(), ".chk") {
		return nil
	}

	checks := cp.Parse(path)
	fmt.Printf("Parsing file %s\n", path)
	for _, check := range checks {
		checkList = append(checkList, check)
	}

	return nil
}

func main() {
	fmt.Println("Loading checks...")

	flag.Parse()
	root := flag.Arg(0)
	err := filepath.Walk(root, visit)
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
		os.Exit(1)
	}

	for _, c := range checkList {
		c.Execute()
		fmt.Println(c)
	}
	fmt.Printf("%d checks executed.\n", len(checkList))
	os.Exit(0)

	// Do jobs with params
	gocron.Every(90).Seconds().Do(taskWithParams, 1, "hello")

	// remove, clear and next_run
	_, time := gocron.NextRun()
	fmt.Println(time)

	// function Start start all the pending jobs
	<-gocron.Start()
}

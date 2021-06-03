package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var origin, destination string
var verbose bool
var pool int

func init() {
	flag.StringVar(&origin, "i", "", "directory of origin")
	flag.StringVar(&destination, "o", "", "directory of destination")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.IntVar(&pool, "p", 1, "number of parallel convertions")
	flag.Parse()

	if origin == "" {
		fmt.Println("directory of origin is missing")
		os.Exit(1)
	}
	if destination == "" {
		destination = origin
	}
}

func findCbrFiles() []string {
	output, err := exec.
		Command("bash", "-c", fmt.Sprintf("find %s -type f -iname *.cbr", origin)).
		Output()
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(output), "\n")
	files := make([]string, 0)
	for _, value := range lines {
		if value == "" {
			continue
		}
		files = append(files, string(value))
	}
	return files
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func cbr2mobi(cbr string) {
	mobi := strings.Replace(cbr, origin, destination, 1)
	mobi = strings.Replace(mobi, ".cbr", ".mobi", 1)
	if fileExists(mobi) {
		if verbose {
			fmt.Printf("File already exists, skipping -- %s\n", mobi)
		}
		return
	}

	if err := os.MkdirAll(filepath.Dir(mobi), 0755); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	cmd := exec.Command("ebook-convert", cbr, mobi)
	if verbose {
		fmt.Println(cmd.String())
	}
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "convertion error: %s\n", err)
	}
}

func clearScreenANSI() {
	fmt.Print("\033[H\033[2J")
}

func printProgress(current, total int) {
	if !verbose {
		clearScreenANSI()
		for i := 0; i < current; i++ {
			fmt.Print(".")
		}
	}
	currentProgress := current * 100 / total
	fmt.Printf("(%d / %d) %d%s\n", current, total, currentProgress, "%")
}

type runner struct {
	sync.Mutex
	sync.WaitGroup
	total, done int
	files       chan string
}

func (r *runner) UpdateProgress(fn func(done int)) {
	r.Lock()
	defer r.Unlock()
	fn(r.done)
	r.done++
}

func (r *runner) NewWorker() {
	defer r.Done()
	for dir := range r.files {
		r.UpdateProgress(func(done int) {
			printProgress(done, r.total)
		})
		cbr2mobi(dir)
	}
}

func newRunner(total int) *runner {
	return &runner{
		WaitGroup: sync.WaitGroup{},
		total:     total,
		files:     make(chan string),
	}
}

func main() {
	files := findCbrFiles()
	runner := newRunner(len(files))
	for i := 0; i < pool; i++ {
		runner.Add(1)
		go runner.NewWorker()
	}
	for _, dir := range files {
		runner.files <- dir
	}
	close(runner.files)
	runner.Wait()
}

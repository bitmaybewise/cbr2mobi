package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var origin, destination string
var verbose bool

func init() {
	flag.StringVar(&origin, "i", "", "directory of origin")
	flag.StringVar(&destination, "o", "", "directory of destination")
	flag.BoolVar(&verbose, "v", false, "verbose output")
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
	return strings.Split(string(output), "\n")
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

func main() {
	files := findCbrFiles()
	total := len(files) - 1
	for i, filename := range files {
		printProgress(i, total)
		if filename == "" {
			continue
		}
		cbr2mobi(filename)
	}
}

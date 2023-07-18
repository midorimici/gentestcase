package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"

	"github.com/midorimici/gentestcase/internal/condition"
	"github.com/midorimici/gentestcase/internal/converter"
	"github.com/midorimici/gentestcase/internal/exporter"
	"github.com/midorimici/gentestcase/internal/filterer"
	"github.com/midorimici/gentestcase/internal/generator"
	"github.com/midorimici/gentestcase/internal/loader"
	"github.com/midorimici/gentestcase/internal/sorter"
)

var (
	inputFilename  = flag.String("input", "cases.yml", "input YAML filename")
	outputFilename = flag.String("output", "data.csv", "output CSV filename")
	isWatching     = flag.Bool("w", false, "watch input file change")
)

func main() {
	flag.Parse()

	// Setup input reader
	var in io.ReadSeeker
	if *inputFilename == "" {
		in = os.Stdin
	} else {
		f, err := os.Open(*inputFilename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		in = f
	}

	if !*isWatching {
		if err := run(in); err != nil {
			log.Fatal(err)
		}
		return
	}

	if *inputFilename == "" {
		log.Fatal("error: cannot watch standard input. You should not use -w flag with empty string -input flag.")
	}

	if err := addWatcher(in); err != nil {
		log.Fatal(err)
	}
}

func addWatcher(in io.ReadSeeker) error {
	const funcName = "addWatcher"

	// Create new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}
	defer watcher.Close()

	// Start listening for events
	go listen(in, watcher)

	// Add a path
	err = watcher.Add(*inputFilename)
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}

	fmt.Printf("Watching %q\n", *inputFilename)

	// Block main goroutine forever
	<-make(chan struct{})

	return nil
}

func listen(in io.ReadSeeker, watcher *fsnotify.Watcher) {
	const funcName = "listen"

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if !event.Has(fsnotify.Write) {
				continue
			}

			// Check if the file is empty
			info, err := os.Stat(*inputFilename)
			if err != nil {
				log.Fatal(fmt.Errorf("%s: %w", funcName, err))
			}

			// If the file is empty, skip running
			if info.Size() == 0 {
				continue
			}

			log.Printf("file modified: %q\n", event.Name)

			if err := rewindFile(in); err != nil {
				log.Fatal(fmt.Errorf("%s: %w", funcName, err))
			}

			if err := run(in); err != nil {
				log.Fatal(fmt.Errorf("%s: %w", funcName, err))
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.Println("error:", err)
		}
	}
}

func rewindFile(in io.ReadSeeker) error {
	const funcName = "rewindFile"

	_, err := in.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}
	return nil
}

func run(in io.Reader) error {
	const funcName = "run"

	fmt.Printf("Read test definitions from %q ...", *inputFilename)

	// Load data from input
	l := loader.New(in)
	d, err := l.Load()
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}

	fmt.Println(" done")

	// Generate all combinations
	g := generator.New(d.Data.Factors)
	cs := g.Generate()

	cslen := len(cs)
	fmt.Printf("Found %d possible combinations\n", cslen)

	// Filter unnecessary cases
	p := condition.NewParser(d.Data)
	f := filterer.New(d.Data.Factors, p, cs)
	fcs, err := f.Filter()
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}

	fcslen := len(fcs)
	fmt.Printf("Eliminated %d combinations, kept %d combinations\n", cslen-fcslen, fcslen)

	// Sort cases
	s := sorter.New(d.Data.Factors, d.OrderedFactors, d.LevelOrders)
	scs := s.Sort(fcs)

	// Convert cases to table
	cv := converter.New(d.Data.Factors, d.OrderedFactors)
	t := cv.ConvertCombinationMapsToTable(scs)

	// Setup output writer
	var out io.Writer
	if *outputFilename == "" {
		out = os.Stdout
	} else {
		f, err := os.Create(*outputFilename)
		if err != nil {
			return fmt.Errorf("%s: %w", funcName, err)
		}
		defer f.Close()
		out = f
	}

	fmt.Printf("Write test cases to %q ...", *outputFilename)

	// Export to CSV
	e := exporter.New(out, d.Data.Factors, d.OrderedFactors)
	if err := e.ExportCSV(t); err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}

	fmt.Println(" done")

	return nil
}

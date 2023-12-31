package runner

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
	"github.com/midorimici/gentestcase/internal/model"
	"github.com/midorimici/gentestcase/internal/sorter"
)

// Runner is used to run the program.
type Runner interface {
	// Run runs the program.
	Run() error
}

type runner struct {
	inputFilename  string
	outputFilename string
	isWatching     bool
	debugFilename  string
}

// New returns a new Runner for given I/O sources and options.
//
// When isWatching is true, the program runs with hot reloading mode.
//
// When debugFilename is not empty, the program outputs debug information.
func New(inputFilename, outputFilename string, isWatching bool, debugFilename string) Runner {
	return &runner{
		inputFilename:  inputFilename,
		outputFilename: outputFilename,
		isWatching:     isWatching,
		debugFilename:  debugFilename,
	}
}

func (r *runner) Run() error {
	flag.Parse()

	// Setup input reader
	var in io.ReadSeeker
	if r.inputFilename == "" {
		in = os.Stdin
	} else {
		f, err := os.Open(r.inputFilename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	if !r.isWatching {
		if err := r.run(in); err != nil {
			return err
		}
		return nil
	}

	if r.inputFilename == "" {
		return fmt.Errorf("error: cannot watch standard input. You should not use -w flag with empty string -input flag")
	}

	if err := r.addWatcher(in); err != nil {
		return err
	}

	return nil
}

func (r *runner) addWatcher(in io.ReadSeeker) error {
	const funcName = "runner.addWatcher"

	// Create new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}
	defer watcher.Close()

	// Start listening for events
	go r.listen(in, watcher)

	// Add a path
	err = watcher.Add(r.inputFilename)
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}

	fmt.Printf("Watching %q\n", r.inputFilename)

	// Block main goroutine forever
	<-make(chan struct{})

	return nil
}

func (r *runner) listen(in io.ReadSeeker, watcher *fsnotify.Watcher) {
	const funcName = "runner.listen"

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
			info, err := os.Stat(r.inputFilename)
			if err != nil {
				log.Fatal(fmt.Errorf("%s: %w", funcName, err))
			}

			// If the file is empty, skip running
			if info.Size() == 0 {
				continue
			}

			fmt.Println()
			log.Printf("file modified: %q\n", event.Name)

			if err := rewindFile(in); err != nil {
				log.Fatal(fmt.Errorf("%s: %w", funcName, err))
			}

			if err := r.run(in); err != nil {
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

func (r *runner) run(in io.Reader) error {
	const funcName = "run"

	fmt.Printf("Read test definitions from %q ...", r.inputFilename)

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
	f := filterer.New(d.Data.Constraints, p, cs, r.debugFilename != "")
	fcs, dcs, err := f.Filter()
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}

	fcslen := len(fcs)
	fmt.Printf("Eliminated %d combinations, kept %d combinations\n", cslen-fcslen, fcslen)

	if err := sortConvertExport(d, fcs, r.outputFilename, false); err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}

	if r.debugFilename != "" {
		if err := sortConvertExport(d, dcs, r.debugFilename, true); err != nil {
			return fmt.Errorf("%s: %w", funcName, err)
		}
	}

	return nil
}

func sortConvertExport(d *loader.Result, combs []model.Combination, outputFilename string, isDebug bool) error {
	const funcName = "sortConvertExport"

	// Sort cases
	s := sorter.New(d.Data.Factors, d.OrderedFactors, d.LevelOrders)
	scs := s.Sort(combs)

	// Convert cases to table
	cv := converter.New(d.Data.Factors, d.OrderedFactors, isDebug)
	t := cv.ConvertCombinationMapsToTable(scs)

	// Setup output writer
	var out io.Writer
	if outputFilename == "" {
		out = os.Stdout
	} else {
		f, err := os.Create(outputFilename)
		if err != nil {
			return fmt.Errorf("%s: %w", funcName, err)
		}
		defer f.Close()
		out = f
	}

	text := "Write test cases to %q ..."
	if isDebug {
		text = fmt.Sprintf("[debug] %s", text)
	}
	fmt.Printf(text, outputFilename)

	// Export to CSV
	e := exporter.New(out)
	if err := e.ExportCSV(t); err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}

	fmt.Println(" done")

	return nil
}

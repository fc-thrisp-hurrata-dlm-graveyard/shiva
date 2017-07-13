package main

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/Laughs-In-Flowers/flip"
	"github.com/Laughs-In-Flowers/log"
	"github.com/Laughs-In-Flowers/shiva/lib/engine"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"

	// initiailize & register providers with graphics package
	_ "github.com/Laughs-In-Flowers/shiva/lib/graphics/providers"
)

type Options struct {
	debug     bool
	formatter string
	provider  string
	dir, file string
}

func defaultOptions() *Options {
	wd, _ := os.Getwd()
	defaultProvider := graphics.DefaultProvider.String()
	return &Options{
		false, "null", defaultProvider, wd, "main.lua",
	}
}

func tFlags(fs *flip.FlagSet, o *Options) *flip.FlagSet {
	fs.BoolVar(&o.debug, "debug", o.debug, "Set engine debug to true.")
	fs.StringVar(&o.formatter, "formatter", o.formatter, "Specify the log formatter.")
	fs.StringVar(&o.provider, "provider", o.provider, "String tag to specify the graphics provder.")
	return fs
}

type Execute func(*Options) error

type Executing interface {
	Run(*Options) error
}

type executing struct {
	has []Execute
}

func NewExecuting(e ...Execute) *executing {
	return &executing{
		e,
	}
}

func (e *executing) Run(o *Options) error {
	for _, v := range e.has {
		err := v(o)
		if err != nil {
			return err
		}
	}
	return nil
}

func debugLog(o *Options) error {
	if o.debug {
		o.formatter = "stdout"
	}
	return nil
}

var topExecute = NewExecuting(debugLog)

func TopCommand(o *Options) flip.Command {
	fs := flip.NewFlagSet("top", flip.ContinueOnError)
	fs = tFlags(fs, o)
	return flip.NewCommand(
		"",
		"shiva",
		"shiva top options",
		1,
		func(c context.Context, a []string) flip.ExitStatus {
			topExecute.Run(o)
			return flip.ExitNo
		},
		fs,
	)
}

func basicErr(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", e)
		os.Exit(-1)
	}
}

func newEngine(o *Options) *engine.Engine {
	configuration := []engine.Config{
		engine.SetLogger(o.formatter),
		engine.SetGraphics(o.provider),
		engine.SetLua(o.dir, o.file),
	}
	v, err := engine.New(o.debug, configuration...)
	if err != nil {
		basicErr(err)
	}
	return v
}

func pFlags(fs *flip.FlagSet, o *Options) *flip.FlagSet {
	fs.StringVar(&o.dir, "dir", o.dir, "The lua directory argument to pass to the engine.")
	fs.StringVar(&o.file, "file", o.file, "The main lua file argument to pass to the engine.")
	return fs
}

func PlayCommand(o *Options) flip.Command {
	fs := flip.NewFlagSet("play", flip.ContinueOnError)
	fs = pFlags(fs, o)
	return flip.NewCommand(
		"",
		"play",
		"shiva play",
		1,
		func(c context.Context, a []string) flip.ExitStatus {
			e := newEngine(o)
			e.Run()
			return flip.ExitSuccess
		},
		fs,
	)
}

var (
	versionPackage string = path.Base(os.Args[0])
	versionTag     string = "No Tag"
	versionHash    string = "No Hash"
	versionDate    string = "No Date"
)

var (
	options     *Options
	currentPath string
	C           *flip.Commander
)

func init() {
	options = defaultOptions()
	log.SetFormatter("shiva_text", log.MakeTextFormatter(versionPackage))
	C = flip.BaseWithVersion(versionPackage, versionTag, versionHash, versionDate)
	C.RegisterGroup("top", 1, TopCommand(options))
	C.RegisterGroup("play", 10, PlayCommand(options))
}

func main() {
	ctx := context.Background()
	C.Execute(ctx, os.Args)
	os.Exit(0)
}

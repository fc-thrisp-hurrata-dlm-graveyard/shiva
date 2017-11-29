package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/Laughs-In-Flowers/flip"
	"github.com/Laughs-In-Flowers/log"
	"github.com/Laughs-In-Flowers/shiva/lib/engine"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"

	// initialize & register providers with graphics package
	_ "github.com/Laughs-In-Flowers/shiva/lib/graphics/providers"
)

type Options struct {
	debug     bool
	formatter string
	provider  string
	file      string
}

func defaultOptions() *Options {
	wd, _ := os.Getwd()
	defaultProvider := graphics.DefaultProvider.String()
	return &Options{
		false, "null", defaultProvider, filepath.Join(wd, "main.lua"),
	}
}

func tFlags(fs *flip.FlagSet, o *Options) *flip.FlagSet {
	fs.BoolVar(&o.debug, "debug", o.debug, "Set engine debug to true.")
	fs.StringVar(&o.file, "file", o.file, "The main lua file argument to pass to the engine.")
	fs.StringVar(&o.formatter, "formatter", o.formatter, "Specify the log formatter.")
	fs.StringVar(&o.provider, "provider", o.provider, "String tag to specify the graphics provder.")
	return fs
}

type Execute func(*Options) error

type executing struct {
	has []Execute
}

func newExecuting(e ...Execute) *executing {
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

var topExecute = newExecuting(debugLog)

func topCommand(o *Options) flip.Command {
	fs := flip.NewFlagSet("top", flip.ContinueOnError)
	fs = tFlags(fs, o)
	return flip.NewCommand(
		"",
		"shiva",
		"shiva top options",
		1,
		false,
		func(c context.Context, a []string) (context.Context, flip.ExitStatus) {
			topExecute.Run(o)
			return c, flip.ExitNo
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
		engine.SetLua(o.file),
	}
	v, err := engine.New(o.debug, configuration...)
	if err != nil {
		basicErr(err)
	}
	return v
}

func playCommand(o *Options) flip.Command {
	return flip.NewCommand(
		"",
		"play",
		"shiva play",
		1,
		false,
		func(c context.Context, a []string) (context.Context, flip.ExitStatus) {
			e := newEngine(o)
			e.Run()
			return c, flip.ExitSuccess
		},
		flip.NewFlagSet("play", flip.ContinueOnError),
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
	F           flip.Flip
)

func init() {
	options = defaultOptions()
	log.SetFormatter("shiva_text", log.MakeTextFormatter(versionPackage))
	F = flip.Base
	F.AddCommand("version", versionPackage, versionTag, versionHash, versionDate).
		AddCommand("help").
		SetGroup("top", -1, topCommand(options)).
		SetGroup("play", 1, playCommand(options))
}

func main() {
	ctx := context.Background()
	F.Execute(ctx, os.Args)
	os.Exit(0)
}

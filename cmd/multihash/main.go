//go:generate curl -sSL https://raw.githubusercontent.com/multiformats/go-multihash/master/version.json -o version.json
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"

	mh "github.com/multiformats/go-multihash"
	mhopts "github.com/multiformats/go-multihash/opts"
	_ "github.com/multiformats/go-multihash/register/all"
	_ "github.com/multiformats/go-multihash/register/miniosha256"
	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  Variables
// ----------------------------------------------------------------------------

// Embed upstream's latest version.json via gogenerate.
//go:embed version.json
var versionJSON []byte

// these variables should be set during build via ldflags.
var (
	version string // release version
	commit  string // git short commit
)

// flag values to hold.
var (
	opts       *mhopts.Options
	checkRaw   string
	checkMh    mh.Multihash
	isQuiet    bool
	isHelp     bool
	isVerLong  bool
	isVerShort bool
)

// help message format.
var usage = `usage: %s [options] [FILE]
Print or check multihash checksums.
With no FILE, or when FILE is -, read standard input.

Options:
`

// exit on error.
var checkErr = func(err error) {
	if err != nil {
		die("error: ", err)
	}
}

// copy of os.Exit to ease testing.
var osExit = os.Exit

// ----------------------------------------------------------------------------
//  Initialization
// ----------------------------------------------------------------------------

// initialize flag options.
func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}

	opts = mhopts.SetupFlags(flag.CommandLine)

	checkStr := "check checksum matches"
	flag.StringVar(&checkRaw, "check", "", checkStr)
	flag.StringVar(&checkRaw, "c", "", checkStr+" (shorthand)")

	helpStr := "display help message"
	flag.BoolVar(&isHelp, "help", false, helpStr)
	flag.BoolVar(&isHelp, "h", false, helpStr+" (shorthand)")

	quietStr := "quiet output (no newline on checksum, no error text)"
	flag.BoolVar(&isQuiet, "quiet", false, quietStr)
	flag.BoolVar(&isQuiet, "q", false, quietStr+" (shorthand)")

	versionStr := "display app version"
	flag.BoolVar(&isVerLong, "version", false, versionStr+" and modules used")
	flag.BoolVar(&isVerShort, "v", false, versionStr+" (shorthand)")
}

// ----------------------------------------------------------------------------
//  Main
// ----------------------------------------------------------------------------

func main() {
	checkErr(parseFlags(opts))

	preRun() // check options such as help and version.

	inp, err := getInput()
	checkErr(err)

	defer inp.Close()

	if checkMh != nil {
		checkErr(opts.Check(inp, checkMh))

		if !isQuiet {
			fmt.Println("OK checksums match (-q for no output)")
		}
	} else {
		checkErr(printHash(opts, inp))
	}
}

// ----------------------------------------------------------------------------
//  Function
// ----------------------------------------------------------------------------

func die(v ...interface{}) {
	if !isQuiet {
		fmt.Fprint(os.Stderr, v...)
		fmt.Fprint(os.Stderr, "\n")
	}

	osExit(1)
}

func getInput() (io.ReadCloser, error) {
	args := flag.Args()

	switch {
	case len(args) < 1:
		return os.Stdin, nil
	case args[0] == "-":
		return os.Stdin, nil
	default:
		f, err := os.Open(args[0])
		if err != nil {
			return nil, fmt.Errorf("failed to open '%s': %s", args[0], err)
		}

		return f, nil
	}
}

func parseFlags(o *mhopts.Options) error {
	flag.Parse()

	if err := o.ParseError(); err != nil {
		return errors.Wrap(err, "failed to parse flags")
	}

	if checkRaw != "" {
		var err error

		checkMh, err = mhopts.Decode(o.Encoding, checkRaw)
		if err != nil {
			return fmt.Errorf("fail to decode multihash '%s': %s", checkRaw, err)
		}
	}

	return nil
}

func preRun() {
	if isHelp {
		flag.Usage()
		osExit(0)
	}

	if isVerLong || isVerShort {
		printVer()
		osExit(0)
	}
}

func printHash(o *mhopts.Options, r io.Reader) error {
	h, err := o.Multihash(r)
	if err != nil {
		return errors.Wrap(err, "failed to calculate its multihash")
	}

	s, err := mhopts.Encode(o.Encoding, h)
	if err != nil {
		return errors.Wrap(err, "failed to encode hash to multihash")
	}

	if isQuiet {
		fmt.Print(s)
	} else {
		fmt.Println(s)
	}

	return nil
}

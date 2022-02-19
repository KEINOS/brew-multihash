package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	mh "github.com/multiformats/go-multihash"
	mhopts "github.com/multiformats/go-multihash/opts"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"
)

// ============================================================================
//  Tests
// ============================================================================

// ----------------------------------------------------------------------------
//  main()
// ----------------------------------------------------------------------------

func Test_main_golden_stdin_quiet(t *testing.T) {
	restore := mockInput(t, []string{
		"-q",
	}, "Hello, world!\n")
	defer restore()

	// quiet option should not add newline
	expect := "QmcwkKyBLujMQitrGSLdtFTzEYSzA7VcfARhFHbe4hZJc4"
	actual := capturer.CaptureOutput(func() {
		main()
	})

	require.Equal(t, expect, actual)
}

func Test_main_golden_stdin_hyphen(t *testing.T) {
	restore := mockInput(t, []string{
		"-",
	}, "Hello, world!\n")
	defer restore()

	// With no "-q" option, it should add newline as well
	expect := "QmcwkKyBLujMQitrGSLdtFTzEYSzA7VcfARhFHbe4hZJc4\n"
	actual := capturer.CaptureOutput(func() {
		main()
	})

	require.Equal(t, expect, actual)
}

func Test_main_golden_file_hash(t *testing.T) {
	// Create test file
	pathFile := filepath.Join(t.TempDir(), t.Name()+"_sample.txt")

	err := os.WriteFile(pathFile, []byte("Hello, world!"), 0o600)
	require.NoError(t, err)

	restore := mockInput(t, []string{pathFile}, "")
	defer restore()

	expect := "QmRfP2G7Nb6SiPZqQxMxtZ1f4hBjY2JGkWvuxvUhkWm6ca\n"
	actual := capturer.CaptureOutput(func() {
		main()
	})

	require.Equal(t, expect, actual)
}

func Test_main_golden_check(t *testing.T) {
	// Create test file
	pathFile := filepath.Join(t.TempDir(), t.Name()+"_sample.txt")

	err := os.WriteFile(pathFile, []byte("Hello, world!"), 0o600)
	require.NoError(t, err)

	// Execute with check option
	restore := mockInput(t, []string{
		"-c", "QmRfP2G7Nb6SiPZqQxMxtZ1f4hBjY2JGkWvuxvUhkWm6ca",
		pathFile,
	}, "")
	defer restore()

	expect := "OK checksums match (-q for no output)\n"
	actual := capturer.CaptureOutput(func() {
		main()
	})

	require.Equal(t, expect, actual)
}

func Test_main_show_help(t *testing.T) {
	restore := mockInput(t, []string{"-h"}, "")
	defer restore()

	oldOsExit := osExit
	defer func() {
		osExit = oldOsExit
	}()

	expectCode := 0

	osExit = func(code int) {
		if code != expectCode {
			t.Logf("help should exit with code %d", expectCode)
			t.Errorf("unexpected exit code %d, expect %d", code, expectCode)
		}

		t.SkipNow()
	}

	_ = capturer.CaptureOutput(func() {
		main()
	})
}

// ----------------------------------------------------------------------------
//  checkErr() (including die())
// ----------------------------------------------------------------------------

func Test_checkErr(t *testing.T) {
	oldOsExit := osExit
	defer func() {
		osExit = oldOsExit
	}()

	osExit = func(code int) {
		if code != 1 {
			t.Fatalf("unexpected exit code %d", code)
		}
	}

	errMsg := "error_" + t.Name()

	out := capturer.CaptureStderr(func() {
		err := errors.New(errMsg)

		checkErr(err)
	})

	assert.Contains(t, out, errMsg)
}

// ----------------------------------------------------------------------------
//  flag.Usage()
// ----------------------------------------------------------------------------

func Test_flagUsage(t *testing.T) {
	out := capturer.CaptureStderr(func() {
		flag.Usage()
	})

	assert.Contains(t, out, "usage:")
	assert.Contains(t, out, "Print or check multihash checksums.")
}

// ----------------------------------------------------------------------------
//  getInput()
// ----------------------------------------------------------------------------

func Test_getInput_file_not_found(t *testing.T) {
	// Dummy file
	pathFile := filepath.Join(t.TempDir(), "unknown.txt")

	restore := mockInput(t, []string{pathFile}, "")
	defer restore()

	flag.Parse()

	r, err := getInput()

	require.Error(t, err, "it should return an error if file not found")
	require.Nil(t, r, "it should be nil on error")
	assert.Contains(t, err.Error(), "failed to open", "it should contain error message")
	assert.Contains(t, err.Error(), pathFile, "the error should contain the file path")
}

// ----------------------------------------------------------------------------
//  parseFlags()
// ----------------------------------------------------------------------------

func Test_parseFlags_wrong_algo(t *testing.T) {
	checkRaw = "QmcwkKyBLujMQitrGSLdtFTzEYSzA7VcfARhFHbe4hZJc4\n"

	restore := mockInput(t, []string{
		"-a", "sha5963", // unknown algorithm
	}, "Hello, world!\n")
	defer restore()

	err := parseFlags(opts)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse flags: algorithm 'sha5963' not")
}

func Test_parseFlags_wrong_hash_vlue(t *testing.T) {
	// Hash sha2-256
	checkRaw = "QmcwkKyBLujMQitrGSLdtFTzEYSzA7VcfARhFHbe4hZJc4\n"

	restore := mockInput(t, []string{
		"-a", "sha3", // unknown algorithm
	}, "Hello, world!\n")
	defer restore()

	oldCheckRaw := checkRaw
	defer func() {
		checkRaw = oldCheckRaw
	}()

	err := parseFlags(opts)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail to decode multihash")
}

// ----------------------------------------------------------------------------
//  preRun()
// ----------------------------------------------------------------------------

func Test_preRun_help(t *testing.T) {
	// Mock osExit
	oldOsExit := osExit
	defer func() {
		osExit = oldOsExit
	}()

	osExit = func(code int) {
		if code != 0 {
			t.Fatalf("unexpected exit code %d", code)
		}
	}

	// Mock isHelp
	oldIsHelp := isHelp
	defer func() {
		isHelp = oldIsHelp
	}()

	isHelp = true // Set show help

	out := capturer.CaptureStderr(func() {
		preRun()
	})

	assert.Contains(t, out, "usage:")
	assert.Contains(t, out, "Print or check multihash checksums.")
}

func Test_preRun_version(t *testing.T) {
	restore := mockInput(t, []string{}, "")
	defer restore()

	// Mock osExit
	oldOsExit := osExit

	defer func() {
		osExit = oldOsExit
	}()

	osExit = func(code int) {
		if code != 0 {
			t.Fatalf("unexpected exit code %d", code)
		}
	}

	isVerLong = true // Show version
	version = "1.2.3"
	commit = "abcdef"

	out := capturer.CaptureStdout(func() {
		preRun()
	})

	assert.Contains(t, out, "v1.2.3-abcdef")
	assert.Contains(t, out, "Modules:")
	assert.Contains(t, out, "go-multihash")
	assert.Contains(t, out, "go-utiles")
}

// ----------------------------------------------------------------------------
//  printHash()
// ----------------------------------------------------------------------------

func Test_printHash_failed_read(t *testing.T) {
	inputDummy := "dummy"
	optsDummy := &mhopts.Options{} // malformed options
	readDummy := bytes.NewReader([]byte(inputDummy))

	err := printHash(optsDummy, readDummy)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to calculate its multihash")
}

func Test_printHash_failed_encode(t *testing.T) {
	inputDummy := "dummy"
	optsDummy := &mhopts.Options{
		Encoding:      "base300", // unknown encoding
		Algorithm:     "sha2-256",
		AlgorithmCode: mh.SHA2_256,
		Length:        len(inputDummy),
	}
	readDummy := bytes.NewReader([]byte("dummy"))

	err := printHash(optsDummy, readDummy)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to encode hash to multihash")
	assert.Contains(t, err.Error(), "unknown encoding: base300")
}

// ----------------------------------------------------------------------------
//  printVer()
// ----------------------------------------------------------------------------

func Test_printVersion(t *testing.T) {
	oldVersionJSON := versionJSON
	oldVersion := version
	oldIsVerShort := isVerShort

	defer func() {
		versionJSON = oldVersionJSON
		version = oldVersion
		isVerShort = oldIsVerShort
	}()

	// Golden
	{
		versionJSON = []byte("{\"version\":\"2.3.4\"}")
		version = ""
		isVerShort = true

		expect := "v2.3.4\n"
		actual := capturer.CaptureStdout(func() {
			printVer()
		})

		assert.Equal(t, expect, actual, "it should use the value of versionJSON")
	}

	// Undefined version
	{
		versionJSON = []byte("")
		version = ""
		isVerShort = true

		expect := "v(undefined)\n"
		actual := capturer.CaptureStdout(func() {
			printVer()
		})

		assert.Equal(t, expect, actual, "it should contain '(undefined)' if version is not defined")
	}
}

// ----------------------------------------------------------------------------
//  uniformVersion()
// ----------------------------------------------------------------------------

func Test_uniformVersion(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect string
	}{
		{"V1.2.3", "v1.2.3"},
		{"1.2.3", "v1.2.3"},
		{"1.2.3-ABCDE", "v1.2.3"},
		{"v0.1.1-0.20211210143450-760ee2c43a7c", "v0.1.1"},
	} {
		expect := test.expect
		actual := uniformVersion(test.input)

		assert.Equal(t, expect, actual)
	}
}

// ============================================================================
//  Helper Functions
// ============================================================================

// Backup variables, mock args and stdin, then return a function to defer restoration
// of the old variables.
func mockInput(t *testing.T, args []string, stdin string) (restore func()) {
	t.Helper()

	// Backup old values
	oldOsArgs := os.Args
	oldOpts := opts
	oldcheckRaw := checkRaw
	oldCheckMh := checkMh
	oldIsQuiet := isQuiet
	oldIsHelp := isHelp
	oldIsVerLong := isVerLong
	oldIsVerShort := isVerShort
	oldCommit := commit

	// Mock os.Args
	os.Args = []string{t.Name()}
	os.Args = append(os.Args, args...)

	// Backup and defer restore os.Stdin
	oldOsStdin := os.Stdin

	// Mock os.Stdin
	alt, closer, err := stdInAlt(t, "Hello, world!\n")
	require.NoError(t, err)

	os.Stdin = alt

	// Restore function
	return func() {
		os.Args = oldOsArgs
		os.Stdin = oldOsStdin
		opts = oldOpts
		checkRaw = oldcheckRaw
		checkMh = oldCheckMh
		isQuiet = oldIsQuiet
		isHelp = oldIsHelp
		isVerLong = oldIsVerLong
		isVerShort = oldIsVerShort
		commit = oldCommit
		_ = closer(alt) // ensure to close
	}
}

// stdInAlt returns Alt(*os.File), closer(func), error
// Alt: intended to replace os.Stdin while test
// closer: Close and remove files internally generated by ALt
func stdInAlt(t *testing.T, s string) (*os.File, func(*os.File) error, error) {
	t.Helper()

	f, err := os.CreateTemp("", "tmp")
	if err != nil {
		return nil, nil, err
	}

	if _, err = f.Write([]byte(s)); err != nil {
		return nil, nil, err
	}

	_, err = f.Seek(0, 0) // jump to head
	require.NoError(t, err)

	closer := func(file *os.File) error {
		if err := file.Close(); err != nil {
			return err
		}

		os.Remove(file.Name())

		return nil
	}

	return f, closer, nil
}

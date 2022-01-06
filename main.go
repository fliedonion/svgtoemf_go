package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

const EXPORT_TYPE = "emf"
const EXPORT_EXT = EXPORT_TYPE

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage : " + os.Args[0] + " filename")
		fmt.Println("    filename : input filename to convert to " + EXPORT_TYPE)

		os.Exit(1)
	}

	// infile := `\\vmware-host\Shared Folders\vmshare\inkscape-emf-load-test\test.fld\test.svg` // OK
	infile := os.Args[1]
	if _, err := os.Stat(infile); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR input file '"+infile+"' missing")
		os.Exit(2)
	}

	const INFILE_PLACEHOLDER = "{{infile}}"
	const OUTFILE_PLACEHOLDER = "{{outfile}}"

	availables := map[string][]string{
		`C:\Program Files\Inkscape\bin\inkscape.exe`: {INFILE_PLACEHOLDER, "--export-filename", OUTFILE_PLACEHOLDER},
		`C:\Program Files\Inkscape\inkscape.exe`:     {INFILE_PLACEHOLDER, "--export-" + EXPORT_TYPE, OUTFILE_PLACEHOLDER},
	}

	fmt.Fprintf(os.Stderr, "Searching Inkscape.\n")
	var foundapp string = ""
	var argPattern []string

	for app, arg := range availables {
		if _, err := os.Stat(app); err == nil {
			foundapp = app
			argPattern = arg
			// fmt.Fprintf(os.Stderr, "FOUND : %s\n", app)
			break
		}
		// fmt.Fprintf(os.Stderr, "MISS : %s\n", app)
	}

	if foundapp == "" {
		fmt.Fprintf(os.Stderr, "inkscape not found!\n")
		os.Exit(1)
	}

	outfile, _ := MakeOutputFile(infile) // second return value is file extension.

	if infile == outfile {
		fmt.Fprintf(os.Stderr, "Terminated. May be input file and output file are same.\n")
		os.Exit(1)
	}

	fmt.Println("convert to '" + outfile + "'")

	var inkscapeArgs []string
	if argPattern != nil {
		inkscapeArgs = make([]string, len(argPattern))
		for i, v := range argPattern {
			v = strings.Replace(v, INFILE_PLACEHOLDER, infile, -1)
			v = strings.Replace(v, OUTFILE_PLACEHOLDER, outfile, -1)
			inkscapeArgs[i] = v
		}
	}

	var cmd *exec.Cmd
	if inkscapeArgs == nil {
		cmd = exec.Command(foundapp)
	} else {
		cmd = exec.Command(foundapp, inkscapeArgs...)
	}

	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Call Inkscape Succeeded")
		fmt.Println("  " + infile)
		fmt.Println("  -> " + outfile)
	}
}

func MakeOutputFile(infile string) (string, string) {
	infileExt := path.Ext(infile)

	if strings.ContainsAny(infileExt, `\/`) {
		// contains path separator, this is not a file extension.
		// like c:\abc.efg\testfile  -> infileExt == .efg\testfile

		infileExt = ""
	}

	outfile := infile + ".emf"
	if infileExt != "" {
		outfile = infile[0:len(infile)-len(infileExt)] + ".emf"
	}

	return outfile, infileExt
}

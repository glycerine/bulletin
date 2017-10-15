package main

import (
	"flag"
	"fmt"
	bull "github.com/glycerine/bulletin"
	"log"
	"os"
	"path"
)

var ProgramName string = path.Base(os.Args[0])

// demonstrate the sequence of calls to DefineFlags() and ValidateConfig()
func main() {

	myflags := flag.NewFlagSet("myflags", flag.ExitOnError)
	cfg := &bull.Config{}
	cfg.DefineFlags(myflags)

	err := myflags.Parse(os.Args[1:])
	err = cfg.ValidateConfig()
	if err != nil {
		log.Fatalf("%s command line flag error: '%s'", ProgramName, err)
	}

	fmt.Printf("flag parsing done, the rest of program goes here...\n")
}

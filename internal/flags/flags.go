package flags

import (
	"flag"
	"fmt"
	"os"
)

var (
	DIR  = flag.String("dir", "./data", "Path to the directory")
	PORT = flag.Int("port", 8080, "Port number")
	HELP = flag.Bool("help", false, "Show the help screen")
)

func HelpShow() {
	fmt.Println(`Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s [-dir <S>] [-port <N>]  
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`)

	wd, _ := os.Getwd()
	fmt.Printf("\nCurrent working directory: %v\n", wd)
}

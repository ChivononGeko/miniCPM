package main

import (
	"flag"
	"fmt"
	"hot-coffee/internal/flags"
	"hot-coffee/internal/router"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	flag.Parse()
	args := os.Args[1:]

	if len(args) != 0 {
		if !(flags.ArgsCheck(args)) {
			fmt.Println("bad input")
			flags.HelpShow()
			os.Exit(1)
		}
	}

	if *flags.HELP {
		flags.HelpShow()
		os.Exit(0)
	}

	absDir, err := filepath.Abs(*flags.DIR)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}

	if len(args) != 0 {
		if !flags.PortDirChecks(*flags.PORT, absDir) {
			fmt.Println("Port value should be in range (1024, 65535] or Invalid directory path")
			os.Exit(1)
		}
	}

	fmt.Printf("The files will be stored at: %s\n", absDir)
	fmt.Printf("Starting server on port %d ...\n", *flags.PORT)

	mux, err := router.SetupRoutes()
	if err != nil {
		slog.Error("Failed to set up routes", "error", err)
		os.Exit(1)
	}

	portStr := ":" + strconv.Itoa(*flags.PORT)
	log.Fatal(http.ListenAndServe(portStr, mux))
}

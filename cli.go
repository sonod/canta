package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		version bool
		run     bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&version, "version", false, "Print version information and quit.")
	flags.BoolVar(&run, "run", false, "running command in event payload")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	if terminal.IsTerminal(0) {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n  %s [OPTIONS] \nOPTIONS\n", os.Args[0], os.Args[0])
		fmt.Fprint(os.Stderr, "  --version  Print version information and quit.\n")
		fmt.Fprint(os.Stderr, "  --run  running command in event payload\n")
	} else {
		p, err := parseEvents()
		if err != nil {
			log.Println("Could not parse event: ", err)
			return ExitCodeError
		}

		if p == "" {
			log.Println("Payload is null: ", err)
			return ExitCodeError
		}

		if run {
			co, err := exec.Command("sh", "-c", p).CombinedOutput()
			if err != nil {
				log.Println("failed: ", p, err)
				return ExitCodeError
			}
			log.Println(string(co))
		}

		log.Println(string(p))
		log.Println("Done")

		return ExitCodeOK
	}
	return ExitCodeOK
}

func parseEvents() (string, error) {
	log.Println("Waiting for events from STDIN...")
	reader := bufio.NewReader(os.Stdin)
	b, err := reader.Peek(1)
	if err != nil {
		return "", err
	}

	if string(b) == "[" {
		log.Println("Reading Consul event")
		ev, err := ParseConsulEvents(reader)
		if err != nil {
			return "", err
		}
		return string(ev.Payload), nil
	} else {
		line, err := reader.ReadString('\n')
		return strings.Trim(line, "\n"), err
	}
}

type ConsulEvent struct {
	ID      string `json:"ID"`
	Name    string `json:"Name"`
	Payload []byte `json:"Payload"`
	LTime   int    `json:"LTime"`
}

type ConsulEvents []ConsulEvent

func ParseConsulEvents(in io.Reader) (*ConsulEvent, error) {
	var evs ConsulEvents
	dec := json.NewDecoder(in)
	if err := dec.Decode(&evs); err != nil {
		return nil, err
	}
	if len(evs) == 0 {
		return nil, fmt.Errorf("No Consul events found")
	}
	ev := &evs[len(evs)-1]
	log.Println("Consul event ID:", ev.ID)
	return ev, nil
}

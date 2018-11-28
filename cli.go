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
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	command, err := parseEvents()
	if err != nil {
		log.Println("Could not parse event: ", err)
		return ExitCodeError
	}

	_, err = exec.Command(command).CombinedOutput()
	if err != nil {
		log.Println("failed: ", command, err)
		return ExitCodeError
	}

	log.Println("Done")

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

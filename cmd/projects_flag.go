package main

import (
	"errors"
	"fmt"
	"strings"
)

// projectsFlag is custom flag type, a slice of compose project names.
type projectsFlag []string

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (i *projectsFlag) String() string {
	return fmt.Sprint(*i)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (i *projectsFlag) Set(value string) error {
	// If we wanted to allow the flag to be set multiple times,
	// accumulating values, we would delete this if statement.
	// That would permit usages such as
	//	-deltaT 10s -deltaT 15s
	// and other combinations.
	if len(*i) > 0 {
		return errors.New("interval flag already set")
	}
	for _, composeProjName := range strings.Split(value, ",") {
		*i = append(*i, composeProjName)
	}
	return nil
}

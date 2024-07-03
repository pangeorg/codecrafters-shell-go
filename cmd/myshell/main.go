package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Command int

const (
	exit Command = iota
	unknown
)

func determine_command(input string) Command {
	match, _ := regexp.MatchString("^(exit)\\s.*\\d", input)
	if match {
		return exit
	}
	return unknown
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		var input, err = bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			os.Stderr.WriteString("reading command")
			continue
		}
		input = strings.Replace(input, "\n", "", -1)
		var command = determine_command(input)

		switch command {
		case exit:
			var split = strings.Split(input, " ")
			if len(split) < 2 {
				var formatted = fmt.Sprintf("Input not understood: '%s'\n", input)
				fmt.Fprint(os.Stdout, formatted)
				continue
			}
			var exit_code, err = strconv.Atoi(split[1])
			if err != nil {
				var formatted = fmt.Sprintf("Input not understood: '%s'\n", input)
				fmt.Fprint(os.Stdout, formatted)
				continue
			}
			os.Exit(exit_code)
		default:
			var formatted = fmt.Sprintf("%s: command not found\n", input)
			fmt.Fprint(os.Stdout, formatted)
		}
	}
}

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
	echo
	unknown
)

var command_map map[Command]string = map[Command]string{exit: "^(exit)\\s.*\\d", echo: "^(echo).*"}

func determine_command(input string) Command {
	var match bool

	for k, v := range command_map {
		match, _ = regexp.MatchString(v, input)
		if match {
			return k
		}
	}
	return unknown
}

func handle_exit(input string) {
	var split = strings.Split(input, " ")
	if len(split) < 2 {
		var formatted = fmt.Sprintf("Input not understood: '%s'\n", input)
		fmt.Fprint(os.Stdout, formatted)
		return
	}
	var exit_code, err = strconv.Atoi(split[1])
	if err != nil {
		var formatted = fmt.Sprintf("Input not understood: '%s'\n", input)
		fmt.Fprint(os.Stdout, formatted)
		return
	}
	os.Exit(exit_code)
}

func handle_echo(input string) {
	var split = strings.Split(input, " ")
	fmt.Fprint(os.Stdout, strings.Join(split[1:], " "))
	fmt.Fprint(os.Stdout, "\n")
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
			handle_exit(input)
		case echo:
			handle_echo(input)
		default:
			var formatted = fmt.Sprintf("%s: command not found\n", input)
			fmt.Fprint(os.Stdout, formatted)
		}
	}
}

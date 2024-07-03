package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Command string

const (
	exit         Command = "exit"
	echo                 = "echo"
	type_builtin         = "type"
	unknown              = "unknown"
)

var command_map map[Command]string = map[Command]string{exit: "^(exit)\\s.*\\d", echo: "^(echo).*", type_builtin: "^(type).*"}

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

func handle_type(input string) {
	var split = strings.Split(input, " ")
	var type_to_check = split[1]
	if _, ok := command_map[Command(type_to_check)]; ok {
		var formatted = fmt.Sprintf("%s is a shell builtin\n", type_to_check)
		fmt.Fprint(os.Stdout, formatted)
		return
	}
	var formatted = fmt.Sprintf("%s: not found\n", type_to_check)
	fmt.Fprint(os.Stdout, formatted)
	return
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
		case type_builtin:
			handle_type(input)
		default:
			var formatted = fmt.Sprintf("%s: command not found\n", input)
			fmt.Fprint(os.Stdout, formatted)
		}
	}
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

func find_in_path(exe string) (string, error) {
	var env_path = os.Getenv("PATH")
	var paths = strings.Split(env_path, ":")
	for _, path := range paths {
		var formatted = fmt.Sprintf("%s/%s", path, exe)
		var matches, err = filepath.Glob(formatted)
		if err == nil && len(matches) > 0 {
			var formatted = fmt.Sprintf("%s is %s\n", exe, matches[0])
			return formatted, nil
		}
	}
	return "", errors.New("executable not found in PATH")
}

func find_in_builtins(exe string) (string, error) {
	if _, ok := command_map[Command(exe)]; ok {
		var formatted = fmt.Sprintf("%s is a shell builtin\n", exe)
		return formatted, nil
	}
	return "", errors.New("executable not found in builtins")
}

func handle_type(input string) {
	var split = strings.Split(input, " ")
	var exe = split[1]
	var formatted string
	var err error

	formatted, err = find_in_builtins(exe)
	if err == nil {
		fmt.Fprint(os.Stdout, formatted)
		return
	}

	formatted, err = find_in_path(exe)
	if err == nil {
		fmt.Fprint(os.Stdout, formatted)
		return
	}
	formatted = fmt.Sprintf("%s: not found\n", exe)
	fmt.Fprint(os.Stdout, formatted)
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

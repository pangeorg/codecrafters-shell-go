package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Builtin string

const (
	exit         Builtin = "exit"
	echo                 = "echo"
	type_builtin         = "type"
	unknown              = "unknown"
)

func handle_external(exe string, args []string) {
	cmd_args := append([]string{filepath.Base(exe)}, args...)
	cmd := exec.Command(exe, cmd_args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s: command not found: %s\n", exe, err)
	}
}

func parse_builtin(exe string) (Builtin, error) {
	switch exe {
	case string(exit):
		return exit, nil
	case string(echo):
		return echo, nil
	case string(type_builtin):
		return type_builtin, nil
	}
	return unknown, errors.New("Not a builtin")
}

func find_executable(exe string) (string, error) {
	var env_path = os.Getenv("PATH")
	var paths = strings.Split(env_path, ":")
	for _, path := range paths {
		var formatted = fmt.Sprintf("%s/%s", path, exe)
		var matches, err = filepath.Glob(formatted)
		if err == nil && len(matches) > 0 {
			return matches[0], nil
		}
	}
	return "", errors.New("executable not found in PATH")
}

func handle_builtin(builtin Builtin, args []string) {
	switch builtin {
	case exit:
		handle_exit(args)
	case echo:
		handle_echo(args)
	case type_builtin:
		handle_type(args)
	default:
		return
	}
}

func handle_type(args []string) {
	var err error
	var exe string

	if len(args) == 0 || len(args) > 1 {
		fmt.Fprint(os.Stdout, "Usage: type 'command'\n")
		return
	}

	_, err = parse_builtin(args[0])
	if err == nil {
		var formatted = fmt.Sprintf("%s is a shell builtin\n", args[0])
		fmt.Fprint(os.Stdout, formatted)
		return
	}

	exe, err = find_executable(args[0])
	if err == nil {
		var formatted = fmt.Sprintf("%s is %s\n", args[0], exe)
		fmt.Fprint(os.Stdout, formatted)
		return
	}
	exe = fmt.Sprintf("%s: not found\n", args[0])
	fmt.Fprint(os.Stdout, exe)
}

func handle_exit(args []string) {
	var exit_code int
	var err error
	if len(args) > 1 {
		var formatted = fmt.Sprintf("Input not understood: exit '%s'\n", args)
		fmt.Fprint(os.Stdout, formatted)
		return
	}
	if len(args) == 0 {
		exit_code = 0
	} else {
		exit_code, err = strconv.Atoi(args[0])
		if err != nil {
			var formatted = fmt.Sprintf("Input not understood: exit '%s'\n", args)
			fmt.Fprint(os.Stdout, formatted)
			return
		}
	}
	os.Exit(exit_code)
}

func handle_echo(args []string) {
	if len(args) > 0 {
		fmt.Fprint(os.Stdout, strings.Join(args, " "))
	}
	fmt.Fprint(os.Stdout, "\n")
}

func main() {
	var input_command string
	var input string
	var err error
	var builtin Builtin
	var exe string

	for {
		var args []string
		fmt.Fprint(os.Stdout, "$ ")
		input, err = bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			os.Stderr.WriteString("reading command")
			continue
		}
		input = strings.Replace(input, "\n", "", -1)
		if input == "" {
			continue
		}
		var cmd_args = strings.Split(input, " ")
		input_command = cmd_args[0]
		if len(cmd_args) > 1 {
			args = cmd_args[1:]
		} else {
			args = []string{}
		}

		builtin, err = parse_builtin(input_command)
		if err == nil {
			handle_builtin(builtin, args)
			continue
		}
		exe, err = find_executable(input_command)
		if err == nil {
			handle_external(exe, args)
			continue
		}
		var formatted = fmt.Sprintf("%s: command not found\n", input)
		fmt.Fprint(os.Stdout, formatted)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		var command, err = bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			os.Stderr.WriteString("reading command")
			continue
		}
		command = strings.Replace(command, "\n", "", -1)
		switch command {
		default:
			var formatted = fmt.Sprintf("%s: command not found\n", command)
			fmt.Fprint(os.Stdout, formatted)
		}
	}
}

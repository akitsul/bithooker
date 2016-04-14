package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/docopt/docopt-go"
	"github.com/seletskiy/hierr"
)

var (
	version = "1.0"
	usage   = `multihooker` + version + `

multihooker is summoned for using multiple hooks in Atlassian Bitbucket
pre-receive git hook.

Multihooker just reads pre-receive hook contents, runs specified program and
gives specified data <stdin> to program as stdin.

You should pass configuration to multihooker as stdin using following syntax:

	<executable> = "<stdin>"

if you want to add multiline stdin data, you should use triple quotes:

	<executable> = """
	<stdin>
	"""

if there is syntax error or any hook exited with non-zero exit status or any
another error occurred, then multihooker will print error to stderr and exit
with exit code 1.

Usage:
	multihooker [options]

Options:
	-h --help		Show this screen.
	--version       Show version.
`
)

func main() {
	_, err := docopt.Parse(usage, nil, true, version, true, true)
	if err != nil {
		panic(err)
	}

	var configuration map[string]string

	metadata, err := toml.DecodeFile("/dev/stdin", &configuration)
	if err != nil {
		fmt.Println()
		fmt.Fprintln(
			os.Stderr,
			hierr.Errorf(err, "multihooker configuration error"),
		)
		os.Exit(2)
	}

	configuration = prepare(configuration)

	for index, executable := range metadata.Keys() {
		stdin, ok := configuration[executable.String()]
		if !ok {
			fmt.Println()
			fmt.Fprintln(
				os.Stderr,
				hierr.Errorf(err, "configuration error", executable),
			)
			os.Exit(2)
		}

		program := exec.Command(executable.String())
		program.Stdin = bytes.NewBufferString(stdin)
		program.Stdout = os.Stdout
		program.Stderr = os.Stderr

		err := program.Run()
		if err != nil {
			fmt.Println()
			fmt.Fprintln(
				os.Stderr,
				hierr.Errorf(err, "hook %s crashed", executable),
			)
			os.Exit(1)
		}

		if index < len(metadata.Keys()) {
			fmt.Println()
			fmt.Println()
		}
	}

	fmt.Println()
}

func prepare(configuration map[string]string) map[string]string {
	prepared := map[string]string{}
	for executable, stdin := range configuration {
		if strings.HasPrefix(stdin, "\n") {
			stdin = strings.TrimPrefix(stdin, "\n")
		}

		if strings.HasSuffix(stdin, "\n\n") {
			stdin = strings.TrimSuffix(stdin, "\n\n") + "\n"
		}

		prepared[executable] = stdin
	}

	return prepared
}

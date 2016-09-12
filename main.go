package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/kovetskiy/bithooks"
	"github.com/reconquest/hierr-go"
)

var (
	version = "1.0"
	usage   = `bithooker ` + version + `

bithooker is summoned for using multiple hooks in Atlassian Bitbucket
pre-receive git hook.

bithooker just reads pre-receive hook contents, runs specified program with
specified args.

You should pass configuration to bithooker as stdin using following syntax:

<hook-name>@<unique-hook-id>
 <args>
 <args>

<another-hook-name>@<another-unique-hook-id>
 <args>
 <args>

* <hook-name> - name of executing hook which will be summoned for accomplishment
      the task. bithooker will call <hook-name> with <args> (one per line),
      pass self stdin to hook, and pass hook stdout/stderr to bitbucket.

* <unique-hook-id> - it's just unique string for your debugging purposes.

* <args> - hook args, should be indented with one space.

If there is syntax error or any hook exited with non-zero exit code or any
another error occurred, then bithooker will print notice to stderr and exit.

Usage:
    bithooker -h | --help
    bithooker --version
    bithooker <hook>...

Options:
    -h --help        Show this screen.
    --version       Show version.
`
)

type output struct {
	newline bool
}

func main() {
	_, err := docopt.Parse(usage, nil, true, version, true, true)
	if err != nil {
		panic(err)
	}

	var (
		hooksDirectory = filepath.Dir(os.Args[0])
		hooksData      = strings.Join(os.Args[1:], "\n")
	)

	stdin, err := ioutil.ReadFile("/dev/stdin")
	if err != nil {
		fatal(err, "can't read stdin")
	}

	hooks, err := bithooks.Decode(hooksData)
	if err != nil {
		fatal(err, "can't decode hooks")
	}

	for _, hook := range hooks {
		stderr := bytes.NewBuffer(nil)
		output := new(output)

		hookExecutable := filepath.Join(hooksDirectory, hook.Name)

		program := exec.Command(hookExecutable, hook.Args...)
		program.Stdin = bytes.NewBuffer(stdin)
		program.Stdout = output
		program.Stderr = stderr
		program.Env = append(
			os.Environ(),
			"HOOK_NAME="+hook.Name,
			"HOOK_ID="+hook.ID,
		)

		err := program.Run()
		if err != nil {
			if !output.newline {
				fmt.Println()
			}
			fatal(
				hierr.Errorf(
					stderr.String(),
					err.Error(),
				),
				"[hook] %s %s", hook.Name, hook.ID,
			)
		}
	}

	fmt.Println()
}

func fatal(err error, msg string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, hierr.Errorf(err, msg, args...).Error())
	os.Exit(1)
}

func (output *output) Write(data []byte) (int, error) {
	if !output.newline {
		fmt.Println()
	}

	return fmt.Fprint(os.Stderr, string(data))
}

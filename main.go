package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/docopt/docopt-go"
	"github.com/kovetskiy/bithooks"
	"github.com/seletskiy/hierr"
)

var (
	version = "1.0"
	usage   = `bithooker` + version + `

bithooker is summoned for using multiple hooks in Atlassian Bitbucket
pre-receive git hook.

bithooker just reads pre-receive hook contents, runs specified program and
gives specified data <stdin> to program as stdin.

You should pass configuration to bithooker as stdin using following syntax:

<hook-name>@<unique-hook-id>
 <data>
 <data>

<another-hook-name>@<another-unique-hook-id>
 <data>
 <data>

* <hook-name> - name of executing hook which will be summoned for accomplishment
      the task. bithooker will call <hook-name>, pass <data> as stdin,
      and print hook stdout/stderr in real-time.

* <unique-hook-id> - it's just unique string for your debugging purposes.

* <data> - hook configuration data, should be indented with one space.

If there is syntax error or any hook exited with non-zero exit code or any
another error occurred, then bithooker will print notice to stderr and exit.

Usage:
	bithooker [options]

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

	rawHooks, err := ioutil.ReadFile("/dev/stdin")
	if err != nil {
		fatal(err, "can't read stdin")
	}

	hooks, err := bithooks.Decode(string(rawHooks))
	if err != nil {
		fatal(err, "can't decode hooks")
	}

	fmt.Println()

	for index, hook := range hooks {
		program := exec.Command(hook.Name)
		program.Stdin = bytes.NewBufferString(hook.Data)
		program.Stdout = os.Stdout
		program.Stderr = os.Stderr
		program.Env = append(
			os.Environ(),
			"HOOK_NAME="+hook.Name,
			"HOOK_ID="+hook.ID,
		)

		err := program.Run()
		if err != nil {
			fatal(err, "hook %s (%s) crashed", hook.ID, hook.Name)
		}

		if index < len(hooks)-1 {
			fmt.Println()
		}
	}

	fmt.Println()
}

func fatal(err error, msg string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, hierr.Errorf(err, msg, args...).Error())
	fmt.Println()
	os.Exit(1)
}

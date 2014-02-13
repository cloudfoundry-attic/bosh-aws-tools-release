package bosh

import (
	"github.com/vito/cmdtest"
	"io"
	"os"
	"os/exec"
)

func Bosh(args ...string) *cmdtest.Session {
	return BoshInDir("", args...)
}

func BoshInDir(dir string, args ...string) *cmdtest.Session {
	cmd := exec.Command("bosh", args...)
	cmd.Dir = dir
	return runCmd(cmd, sessionStarterSeparateOutput)
}

func BoshCombinedOutput(args ...string) *cmdtest.Session {
	cmd := exec.Command("bosh", args...)
	return runCmd(cmd, sessionStarterCombinedOutput)
}

var sessionStarterSeparateOutput = func(cmd *exec.Cmd) (*cmdtest.Session, error) {
	teeStdout := func(out io.Writer) io.Writer {
		return verboseOutputWriter(out, os.Stdout)
	}
	teeStderr := func(out io.Writer) io.Writer {
		return verboseOutputWriter(out, os.Stderr)
	}
	return cmdtest.StartWrapped(cmd, teeStdout, teeStderr)
}

var sessionStarterCombinedOutput = func(cmd *exec.Cmd) (*cmdtest.Session, error) {
	var stdoutIn io.Writer

	teeStdout := func(out io.Writer) io.Writer {
		stdoutIn = out
		return verboseOutputWriter(out, os.Stdout)
	}
	teeStderr := func(_ io.Writer) io.Writer {
		if stdoutIn == nil {
			panic("stdout must be wrapped first")
		}
		return verboseOutputWriter(stdoutIn, os.Stderr)
	}
	return cmdtest.StartWrapped(cmd, teeStdout, teeStderr)
}

type sessionStarterType func(cmd *exec.Cmd) (*cmdtest.Session, error)

func runCmd(cmd *exec.Cmd, sessionStarter sessionStarterType) *cmdtest.Session {
	sess, err := sessionStarter(cmd)
	if err != nil {
		panic(err)
	}

	return sess
}

func verboseOutputWriter(out, secondary io.Writer) io.Writer {
	if verboseOutputEnabled() {
		return io.MultiWriter(out, secondary)
	}
	return out
}

func verboseOutputEnabled() bool {
	verbose := os.Getenv("VERBOSE_OUTPUT")
	return verbose == "yes" || verbose == "true"
}

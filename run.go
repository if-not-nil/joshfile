package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type ErrorLog struct {
	Err    error
	Stdout string
	Stderr string
	Cmd    string
}

func Run(task Task) []ErrorLog {
	var errors []ErrorLog
	var mu sync.Mutex

	if !task.Async {
		for _, cmd := range task.Cmds {
			var start time.Time
			if task.Report {
				start = time.Now()
			}

			errLog := execute(cmd, task.Silent, task.Direct)
			if errLog.Err != nil {
				errors = append(errors, errLog)
			}

			if task.Report {
				duration := time.Since(start)
				PrintLog(cmd, "done in", duration)
			}
		}
	} else {
		var wg sync.WaitGroup
		for _, cmd := range task.Cmds {
			wg.Add(1)
			go func(cmd string) {
				defer wg.Done()

				var start time.Time
				if task.Report {
					start = time.Now()
				}

				errLog := execute(cmd, task.Silent, task.Direct)
				if errLog.Err != nil {
					mu.Lock()
					errors = append(errors, errLog)
					mu.Unlock()
				}

				if task.Report {
					duration := time.Since(start)
					PrintLog(cmd, ": done in", duration)
				}
			}(cmd)
		}
		wg.Wait()
	}

	return errors
}

func execute(cmd string, silent bool, direct bool) ErrorLog {
	if cmd == "" {
		return ErrorLog{}
	}

	var command *exec.Cmd
	if !direct {
		command = exec.Command("sh", "-c", cmd)
	} else {
		cmdArr := strings.Fields(cmd)
		name := cmdArr[0]
		args := cmdArr[1:]
		command = exec.Command(name, args...)
	}

	command.Env = os.Environ()

	var stdoutBuf, stderrBuf bytes.Buffer
	if silent {
		command.Stdout = &stdoutBuf
		command.Stderr = &stderrBuf
	} else {
		command.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		command.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	}

	err := command.Run()
	if err != nil {
		return reportCommandError(cmd, err, stdoutBuf.String(), stderrBuf.String())
	}

	return ErrorLog{}
}

// reportCommandError structures command outputs into an ErrorLog to make it easier to modify in the future
func reportCommandError(cmdStr string, err error, stdout, stderr string) ErrorLog {
	return ErrorLog{
		Err:    err,
		Stdout: stdout,
		Stderr: stderr,
		Cmd:    cmdStr,
	}
}

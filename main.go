package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/goccy/go-yaml"
)

func main() {
	arg := getArgs()

	switch arg {
	case "--init":
		initFile()
		return
	case "--help", "-h":
		help()
		return
	case "--man":
		showMan()
		return
	case "--list":
		config := parseFile()
		listTasks(config)
		return
	}

	task := getTask(arg)
	if task == nil {
		PrintError("invalid task or task does not exist. run 'josh --list' to see available tasks.")
		return
	}

	errs := Run(*task)
	for _, err := range errs {
		printErrorDetails(err)
	}
}

func getArgs() string {
	if len(os.Args) < 2 {
		return ""
	}
	return os.Args[1]
}

func parseFile() Config {
	data, err := os.ReadFile("./josh.yaml")
	if err != nil {
		log.Fatalf("Error reading josh.yaml: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Error unmarshaling josh.yaml: %v", err)
	}
	return config
}

func getTask(arg string) *Task {
	config := parseFile()
	if arg == "" {
		for _, task := range config.Tasks {
			return &task
		}
	} else {
		if task, exists := config.Tasks[arg]; exists {
			return &task
		}
	}
	return nil
}

func listTasks(config Config) {
	for name, task := range config.Tasks {
		var head string
		if task.Silent {
			head = "[silent] " + head
		}
		if task.Async {
			head = "[async] " + head
		}
		PrintHead(head, name)
		for _, cmd := range task.Cmds {
			fmt.Println("  ", cmd)
		}
	}
}

//go:embed sample/josh.yaml
var sample []byte

func initFile() {
	path := "josh.yaml"
	if _, err := os.ReadFile(path); err == nil {
		PrintError("joshfile already exists")
		return
	}

	err := os.WriteFile(path, sample, 0644)
	if err != nil {
		PrintError("couldn't write josh.yaml:", err)
		return
	}
	PrintLog("joshfile created at", path)
	PrintLog("try the commands from 'josh --list' to see available tasks")
}

//go:embed sample/help.txt
var helpFile []byte

func help() {
	fmt.Println(string(helpFile))
}

//go:embed sample/josh.1
var manPage []byte

func showMan() {
	tmpfile, err := os.CreateTemp("", "josh*.1")
	if err != nil {
		log.Println("error creating temp file:", err)
		fmt.Print(string(manPage)) // fallback
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(manPage); err != nil {
		log.Println("error writing to temp file:", err)
		return
	}
	tmpfile.Close()

	cmd := exec.Command("man", tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Println("error running man:", err)
		fmt.Print(string(manPage)) // fallback
	}
}

func printErrorDetails(err ErrorLog) {
	PrintErrHead("==========")
	PrintError(err.Err, "in", err.Cmd)
	fmt.Println(err.Stderr)
	if err.Stdout != "" {
		PrintLog(err.Stdout)
	}
	PrintErrHead("==========")
}

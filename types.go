package main

type Task struct {
	Silent bool     `yaml:"silent,omitempty"`
	Cmds   []string `yaml:"cmds,omitempty"`
	Async  bool     `yaml:"async,omitempty"`
	Report bool     `yaml:"report,omitempty"`
	Direct bool     `yaml:"direct,omitempty"`
}

type Config struct {
	Tasks map[string]Task `yaml:"tasks"`
}

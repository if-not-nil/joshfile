package main

import "fmt"

const (
	ColorReset      = "\033[0m"
	ColorRed        = "\033[31m"
	ColorCyan       = "\033[36m"
	ColorBoldYellow = "\033[1;33m"
)

func PrintError(args ...any) {
	fmt.Print(ColorRed, "[err] ", ColorReset, fmt.Sprintln(args...))
}

func PrintLog(args ...any) {
	fmt.Print(ColorCyan, "[log] ", ColorReset, fmt.Sprintln(args...))
}

func PrintHead(args ...any) {
	fmt.Print(ColorBoldYellow, fmt.Sprintln(args...), ColorReset)
}

func PrintErrHead(args ...any) {
	fmt.Print("\033[1;31m", fmt.Sprintln(args...), ColorReset)
}


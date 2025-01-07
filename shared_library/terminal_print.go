package main

import "C"
import (
	"fmt"
	"log"
	"strings"
)

const styleReset string = "\033[0m"

// Map for styles
var styleMap = map[string]string{
	// Colors (foreground)
	"[style.yellow]":         "\033[33m",
	"[style.red]":            "\033[31m",
	"[style.green]":          "\033[32m",
	"[style.blue]":           "\033[34m",
	"[style.magenta]":        "\033[35m",
	"[style.cyan]":           "\033[36m",
	"[style.white]":          "\033[37m",
	"[style.black]":          "\033[30m",
	"[style.gray]":           "\033[90m",
	"[style.bright.red]":     "\033[91m",
	"[style.bright.green]":   "\033[92m",
	"[style.bright.yellow]":  "\033[93m",
	"[style.bright.blue]":    "\033[94m",
	"[style.bright.magenta]": "\033[95m",
	"[style.bright.cyan]":    "\033[96m",
	"[style.bright.white]":   "\033[97m",

	// Font styles
	"[style.bold]":          "\033[1m",
	"[style.italic]":        "\033[3m",
	"[style.underline]":     "\033[4m",
	"[style.strikethrough]": "\033[9m",

	// General
	"[style.reset]": styleReset,
}

func logPrint(args ...interface{}) {
	// logPrint simply logs depending to `DEBUG` variable.
	//
	// Parameters:
	//   - args: A `varg` list to intercept logging arguments and pass them to log function.

	if DEBUG {
		out := ""
		for _, arg := range args {
			out = fmt.Sprintf("%s%s", out, arg)
		}
		log.Print(out)
	}
}

func formatStyles(arg string) string {
	// formatStyles search and replace all dedicated style-rectangular-braced tags into style escape codes.
	//
	// Parameters:
	//   - arg: An `arg` string which is going to be searched for style tags.

	out := arg
	if strings.Contains(arg, "[style.") {
		for style, code := range styleMap {
			out = strings.ReplaceAll(out, style, code)
			if !strings.Contains(out, "[style.") {
				return out // speed up a bit
			}
		}
	}
	return out
}

func print(args ...interface{}) {
	// print simply text on terminal.
	//
	// Parameters:
	//   - args: A `varg` list to intercept fmt.Print arguments and pass them to print function.

	out := ""
	for _, arg := range args {
		argString := fmt.Sprintf("%s", arg)
		out = fmt.Sprintf("%s%s", out, formatStyles(argString))
	}
	fmt.Print(out, styleReset)
}

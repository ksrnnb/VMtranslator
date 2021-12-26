package main

import (
	"os"

	"github.com/ksrnnb/VMtranslator/parser"
)

func main() {
	f, err := os.Open("add.vm")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	p := parser.NewParser(f)

	for {
		if !p.HasMoreCommands() {
			break
		}

		p.Advance()
	}
}

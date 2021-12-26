package parser

import (
	"bufio"
	"io"
)

type Parser struct {
	input          io.Reader
	currentCommand string
	isDone         bool
	scanner        *bufio.Scanner
}

const (
	C_ARITHMETIC = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCITON
	C_RETURN
	C_CALL
)

func NewParser(input io.Reader) *Parser {
	scanner := bufio.NewScanner(input)
	return &Parser{
		input:          input,
		currentCommand: "",
		isDone:         false,
		scanner:        scanner,
	}
}

func (p Parser) HasMoreCommands() bool {
	return p.isDone
}

func (p *Parser) Advance() {
	if p.isDone {
		return
	}

	if !p.scanner.Scan() {
		p.isDone = true
		return
	}

	p.currentCommand = p.scanner.Text()
}

func (p Parser) CommandType() int {
	// TODO: implements...
	return 0
}

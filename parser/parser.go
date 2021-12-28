package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
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
	C_FUNCTION
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
	return !p.isDone
}

func (p *Parser) Advance() {
	if p.isDone {
		return
	}

	if !p.scanner.Scan() {
		p.isDone = true
		return
	}

	cmd := p.scanner.Text()

	trimmedCmd := strings.Trim(cmd, " ")

	if len(trimmedCmd) == 0 {
		p.Advance()
		return
	}

	if trimmedCmd[0] == '/' && trimmedCmd[1] == '/' {
		p.Advance()
		return
	}

	p.currentCommand = p.scanner.Text()
}

func (p Parser) CommandType() (int, error) {
	commands := strings.Split(p.currentCommand, " ")

	switch commands[0] {
	case "push":
		return C_PUSH, nil
	case "pop":
		return C_POP, nil
	case "label":
		return C_LABEL, nil
	case "goto":
		return C_GOTO, nil
	case "if-goto":
		return C_IF, nil
	case "function":
		return C_FUNCTION, nil
	case "call":
		return C_CALL, nil
	case "return":
		return C_RETURN, nil
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return C_ARITHMETIC, nil
	}

	return 0, errors.New("command type is invalid")
}

// 今のコマンドの最初の引数を返す
func (p Parser) Arg1() (string, error) {
	cmdType, err := p.CommandType()

	if err != nil {
		return "", err
	}

	commands := strings.Split(p.currentCommand, " ")

	if cmdType == C_ARITHMETIC {
		return commands[0], nil
	}

	return commands[1], nil
}

// 今のコマンドの2番目の引数を返す
func (p Parser) Arg2() (int, error) {
	cmdType, err := p.CommandType()

	if err != nil {
		return 0, err
	}

	if !(cmdType == C_PUSH ||
		cmdType == C_POP ||
		cmdType == C_FUNCTION ||
		cmdType == C_CALL) {
		return 0, fmt.Errorf("arg2 cannot be called in command type %v", cmdType)
	}

	commands := strings.Split(p.currentCommand, " ")
	return strconv.Atoi(commands[2])
}

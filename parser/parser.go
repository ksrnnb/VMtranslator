package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ksrnnb/VMtranslator/command"
)

type Parser struct {
	input          io.Reader
	currentCommand string
	isDone         bool
	scanner        *bufio.Scanner
}

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
		return command.C_PUSH, nil
	case "pop":
		return command.C_POP, nil
	case "label":
		return command.C_LABEL, nil
	case "goto":
		return command.C_GOTO, nil
	case "if-goto":
		return command.C_IF, nil
	case "function":
		return command.C_FUNCTION, nil
	case "call":
		return command.C_CALL, nil
	case "return":
		return command.C_RETURN, nil
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return command.C_ARITHMETIC, nil
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

	if cmdType == command.C_ARITHMETIC {
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

	if !(cmdType == command.C_PUSH ||
		cmdType == command.C_POP ||
		cmdType == command.C_FUNCTION ||
		cmdType == command.C_CALL) {
		return 0, fmt.Errorf("arg2 cannot be called in command type %v", cmdType)
	}

	commands := strings.Split(p.currentCommand, " ")
	trimmed := strings.Split(commands[2], "\t")
	return strconv.Atoi(trimmed[0])
}

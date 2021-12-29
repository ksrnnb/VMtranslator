package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ksrnnb/VMtranslator/command"
	"github.com/ksrnnb/VMtranslator/parser"
	"github.com/ksrnnb/VMtranslator/writer"
)

func main() {
	err := run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	args := os.Args

	if len(args) != 2 {
		return fmt.Errorf("%s: 1 argument should be given", args[0])
	}

	fileInfo, err := os.Stat(args[1])

	if err != nil {
		return fmt.Errorf("%s: cannot read file or directory: %s", args[0], args[1])
	}

	outFileName := strings.Split(fileInfo.Name(), ".")[0] + ".asm"
	file, err := os.Create(outFileName)

	if err != nil {
		return fmt.Errorf("%s: cannot create new file: %s", args[0], outFileName)
	}

	defer file.Close()

	cw := writer.NewCodeWriter(file)

	if fileInfo.IsDir() {
		dirFunc(args[1], cw)
	} else {
		if !isVmFile(args[1]) {
			return fmt.Errorf("%s: file is not vm file: %s", args[0], args[1])
		}

		handleFile(args[1], cw)
	}

	return nil
}

func isVmFile(path string) bool {
	return filepath.Ext(path) == ".vm"
}

func dirFunc(root string, cw *writer.CodeWriter) error {
	err := filepath.Walk(root,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("walk function cannot start: %v", err)
			}

			if info.IsDir() {
				return nil
			}

			if !isVmFile(path) {
				return nil
			}

			return handleFile(path, cw)
		})

	if err != nil {
		return fmt.Errorf("walk function error: %v", err)
	}

	return nil
}

func handleFile(path string, cw *writer.CodeWriter) error {
	f, err := os.Open(path)

	if err != nil {
		return err
	}

	parser := parser.NewParser(f)

	for {
		if !parser.HasMoreCommands() {
			break
		}

		parser.Advance()

		cmdType, err := parser.CommandType()

		if err != nil {
			return fmt.Errorf("handle file: %v", err)
		}

		if cmdType == command.C_PUSH || cmdType == command.C_POP {
			arg1, err := parser.Arg1()

			if err != nil {
				return fmt.Errorf("handle file: %v", err)
			}

			arg2, err := parser.Arg2()

			if err != nil {
				return fmt.Errorf("handle file: %v", err)
			}

			cw.WritePushPop(cmdType, arg1, arg2)
		}
	}

	return err
}

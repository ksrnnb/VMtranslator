package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ksrnnb/VMtranslator/writer"
)

func main() {
	err := run()

	if err != nil {
		fmt.Println(err)
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

	fmt.Println("fileName: ", fileInfo.Name())

	outFileName := strings.Split(fileInfo.Name(), ".")[0] + ".asm"
	file, err := os.Create(outFileName)

	if err != nil {
		return fmt.Errorf("%s: cannot create new file: %s", args[0], outFileName)
	}

	defer file.Close()

	writer.NewCodeWriter(file)

	// p := parser.NewParser(f)

	// for {
	// 	if !p.HasMoreCommands() {
	// 		break
	// 	}

	// 	p.Advance()
	// }

	return nil
}

package writer

import (
	"fmt"
	"io"

	"github.com/ksrnnb/VMtranslator/command"
)

type CodeWriter struct {
	out      io.Writer
	fileName string
}

func NewCodeWriter(out io.Writer) *CodeWriter {
	return &CodeWriter{out: out}
}

func (cw *CodeWriter) SetFileName(fileName string) {
	cw.fileName = fileName
}

func (cw CodeWriter) WriteArithmetic(cmd string) {
	// do something
}

func (cw CodeWriter) WritePushPop(cmd int, segment string, index int) {
	if cmd == command.C_PUSH {
		cw.writePush(segment, index)
	} else if cmd == command.C_POP {
		cw.writePop(segment, index)
	}
}

func (cw *CodeWriter) writePush(segment string, index int) {
	cw.write(fmt.Sprintf("@%d\nD=A\n", index))
}

func (cw *CodeWriter) writePop(segment string, index int) {
	// do something
}

func (cw CodeWriter) write(str string) error {
	_, err := cw.out.Write([]byte(str))

	return err
}

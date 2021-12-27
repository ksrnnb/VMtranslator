package writer

import "io"

type CodeWriter struct {
	out io.Writer
}

func NewCodeWriter(out io.Writer) *CodeWriter {
	return &CodeWriter{out: out}
}

func (cw CodeWriter) SetFileName(fileName string) {
	// do something
}

func (cw CodeWriter) WriteArithmetic(command string) {
	// do something
}

func (cw CodeWriter) WritePushPop(command int, segment string, index int) {
	// do something
}

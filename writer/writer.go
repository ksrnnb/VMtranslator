package writer

import "io"

type CodeWriter struct {
	out io.WriteCloser
}

func NewCodeWriter(out io.WriteCloser) *CodeWriter {
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

func (cw CodeWriter) Close() {
	cw.out.Close()
}

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
	cw.write([]string{
		fmt.Sprintf("@%d", index),
		"D=A",
	})

	cw.write([]string{
		"@SP",   // SPは0
		"A=M",   // AレジスタにRAM[0]を代入、すなわちSPが指すアドレスをAに代入 -> 次からMのアドレスをSPの指す位置に変更
		"M=D",   // Dの値をstackにpushする
		"@SP",   // Mのアドレスを0に戻す
		"M=M+1", // SPが指すアドレスをインクリメントする
	})
}

func (cw *CodeWriter) writePop(segment string, index int) {
	// do something
}

func (cw CodeWriter) write(strs []string) error {
	for _, str := range strs {
		_, err := cw.out.Write([]byte(str + "\n"))

		if err != nil {
			return err
		}
	}

	return nil
}

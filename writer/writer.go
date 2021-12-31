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
	switch cmd {
	case "add":
		cw.writeAdd()
	}
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

	cw.writePushDRegister()
}

func (cw *CodeWriter) writePop(segment string, index int) {

}

func (cw *CodeWriter) writeAdd() {
	cw.write([]string{
		"@SP",   // はじめのRAM[0]=r0とする
		"M=M-1", // SPをデクリメントしてデータが入っている一番上の位置に移動(r0-1)
		"A=M",   // A = RAM[r0-1] => M = RAM[r0-1]
		"D=M",   // D = RAM[r0-1]
		"@SP",   // A = 0, M = RAM[0] = r0-1
		"M=M-1", // M = r0-2
		"A=M",   // A = r0-2
		"D=D+M", // D = RAM[r0-1] + RAM[r0-2]
	})

	cw.writePushDRegister()
}

// Dレジスタの値をpushする
func (cw *CodeWriter) writePushDRegister() {
	cw.write([]string{
		"@SP",   // SPは0
		"A=M",   // AレジスタにRAM[0]を代入、すなわちSPが指すアドレスをAに代入 -> 次からMのアドレスをSPの指す位置に変更
		"M=D",   // Dの値をstackにpushする
		"@SP",   // Mのアドレスを0に戻す
		"M=M+1", // SPが指すアドレスをインクリメントする
	})
}

func (cw *CodeWriter) write(strs []string) error {
	for _, str := range strs {
		_, err := cw.out.Write([]byte(str + "\n"))

		if err != nil {
			return err
		}
	}

	return nil
}

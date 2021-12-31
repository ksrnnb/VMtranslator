package writer

import (
	"fmt"
	"io"

	"github.com/ksrnnb/VMtranslator/command"
)

type CodeWriter struct {
	out        io.Writer
	fileName   string
	compNumber int
}

func NewCodeWriter(out io.Writer) *CodeWriter {
	return &CodeWriter{out: out, compNumber: 0}
}

func (cw *CodeWriter) SetFileName(fileName string) {
	cw.fileName = fileName
}

func (cw *CodeWriter) WriteArithmetic(cmd string) {
	switch cmd {
	case "add":
		cw.writeCommonArithmetic()
		cw.write([]string{"D=M+D"}) // M+D => stackなので右辺のDは後にpushした値。Mは先にpushした値
		cw.writePushDRegister()
	case "sub":
		cw.writeCommonArithmetic()
		cw.write([]string{"D=M-D"})
		cw.writePushDRegister()
	case "neg":
		cw.writeNeg()
	case "eq", "gt", "lt":
		cw.writeCommonArithmetic()
		cw.writeComparison(cmd)
		cw.writePushDRegister()
	case "and":
		cw.writeCommonArithmetic()
		cw.write([]string{"D=M&D"})
		cw.writePushDRegister()
	case "or":
		cw.writeCommonArithmetic()
		cw.write([]string{"D=M|D"})
		cw.writePushDRegister()
	case "not":
		cw.writeNot()
	}
}

func (cw *CodeWriter) WritePushPop(cmd int, segment string, index int) {
	if cmd == command.C_PUSH {
		cw.writePush(segment, index)
	} else if cmd == command.C_POP {
		cw.writePop(segment, index)
	}
}

func (cw *CodeWriter) writePush(segment string, index int) {
	switch segment {
	case "constant":
		cw.writePushConstant(index)
	case "local", "argument", "this", "that":
		cw.writePushSegment(segment, index)
	case "temp":
		cw.writePushTempSegment(index)
	}
}

func (cw *CodeWriter) writePop(segment string, index int) {
	switch segment {

	case "local", "argument", "this", "that":
		cw.writePopSegment(segment, index)
	case "temp":
		cw.writePopTempSegment(index)
	}
}

func (cw *CodeWriter) writeCommonArithmetic() {
	cw.write([]string{
		"@SP",   // はじめのRAM[0]=r0とする
		"M=M-1", // SPをデクリメントしてデータが入っている一番上の位置に移動(r0-1)
		"A=M",   // A = RAM[r0-1] => M = RAM[r0-1]
		"D=M",   // D = RAM[r0-1]
		"@SP",   // A = 0, M = RAM[0] = r0-1
		"M=M-1", // M = r0-2
		"A=M",   // A = r0-2
	})
}

func (cw *CodeWriter) writeComparison(cmd string) {
	var assemblyCmd string

	switch cmd {
	case "eq":
		assemblyCmd = "JEQ"
	case "gt":
		assemblyCmd = "JGT"
	case "lt":
		assemblyCmd = "JLT"
	}

	num := cw.getCompLabelNumber()

	cw.write([]string{
		"D=M-D", // M-D => stackなので右辺のDは後にpushした値。Mは先にpushした値
		fmt.Sprintf("@comp.%d.true", num),
		fmt.Sprintf("D;%s", assemblyCmd),
		fmt.Sprintf("@comp.%d.false", num), // 上の式でjumpしない場合 => false
		"0;JMP",
		fmt.Sprintf("(comp.%d.true)", num),
		"D=-1",
		fmt.Sprintf("@comp.%d.fin", num),
		"0;JMP",
		fmt.Sprintf("(comp.%d.false)", num),
		"D=0",
		fmt.Sprintf("(comp.%d.fin)", num),
	})
}

func (cw *CodeWriter) writeNeg() {
	cw.write([]string{
		"@SP",
		"A=M-1",
		"M=-M",
	})
}

func (cw *CodeWriter) writeNot() {
	cw.write([]string{
		"@SP",
		"A=M-1",
		"M=!M",
	})
}

func (cw *CodeWriter) writePushConstant(index int) {
	cw.write([]string{
		fmt.Sprintf("@%d", index),
		"D=A",
	})
	cw.writePushDRegister()
}

// segmentのindex番地のアドレスの値を、SPにpushする
func (cw *CodeWriter) writePushSegment(segment string, index int) {
	var seg string
	switch segment {
	case "local":
		seg = "LCL"
	case "argument":
		seg = "ARG"
	case "this":
		seg = "THIS"
	case "that":
		seg = "THAT"
	}

	cw.write([]string{
		fmt.Sprintf("@%d", index),
		"D=A",
		fmt.Sprintf("@%s", seg),
		"A=M+D",
		"D=M",
	})

	cw.writePushDRegister()
}

// temp(RAM[5] ~ RAM[12]のindex番地のアドレスの値を、SPにpushする
func (cw *CodeWriter) writePushTempSegment(index int) {
	tempIndex := 5 + index

	cw.write([]string{
		fmt.Sprintf("@%d", tempIndex),
		"D=M",
	})

	cw.writePushDRegister()
}

// SPの値をpopして、segmentのindex番地のアドレスに代入する
func (cw *CodeWriter) writePopSegment(segment string, index int) {
	var seg string
	switch segment {
	case "local":
		seg = "LCL"
	case "argument":
		seg = "ARG"
	case "this":
		seg = "THIS"
	case "that":
		seg = "THAT"
	}

	cw.write([]string{
		"@SP",   // はじめのRAM[0]=r0とする
		"M=M-1", // SPをデクリメントしてデータが入っている一番上の位置に移動(r0-1)
		"A=M",   // A = RAM[r0-1] => M = RAM[r0-1]
		"D=M",   // D = RAM[r0-1]
		fmt.Sprintf("@%s", seg),
		"A=M",
	})

	for i := 0; i < index; i++ {
		cw.write([]string{
			"A=A+1",
		})
	}

	cw.write([]string{
		"M=D",
	})
}

// SPからpopした値をtempに代入
func (cw *CodeWriter) writePopTempSegment(index int) {
	tempIndex := 5 + index

	cw.write([]string{
		"@SP",   // はじめのRAM[0]=r0とする
		"M=M-1", // SPをデクリメントしてデータが入っている一番上の位置に移動(r0-1)
		"A=M",   // A = RAM[r0-1] => M = RAM[r0-1]
		"D=M",   // D = RAM[r0-1]
		fmt.Sprintf("@%d", tempIndex),
		"M=D",
	})
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

func (cw *CodeWriter) getCompLabelNumber() int {
	cw.compNumber++
	return cw.compNumber
}

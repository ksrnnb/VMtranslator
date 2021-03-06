package writer

import (
	"fmt"
	"io"
	"strings"

	"github.com/ksrnnb/VMtranslator/command"
)

type CodeWriter struct {
	out                 io.Writer
	fileName            string
	compNumber          int
	returnNumber        int
	currentFunctionName string
}

func NewCodeWriter(out io.Writer) *CodeWriter {
	return &CodeWriter{out: out, compNumber: 0, returnNumber: 0}
}

func (cw *CodeWriter) SetFileName(fileName string) {
	names := strings.Split(fileName, "/")
	name := names[len(names)-1]
	cw.fileName = name
}

// VM初期化
func (cw *CodeWriter) WriteInit() {
	cw.write([]string{
		"@256",
		"D=A",
		"@SP",
		"M=D",
	})

	cw.WriteCall("Sys.init", 0)
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
	case "pointer":
		cw.writePushPointerSegment(index)
	case "static":
		cw.writePushStaticSegment(index)
	}
}

func (cw *CodeWriter) writePop(segment string, index int) {
	switch segment {

	case "local", "argument", "this", "that":
		cw.writePopSegment(segment, index)
	case "temp":
		cw.writePopTempSegment(index)
	case "pointer":
		cw.writePopPointerSegment(index)
	case "static":
		cw.writePopStaticSegment(index)
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

func (cw *CodeWriter) writePushPointerSegment(index int) {
	switch index {
	case 0:
		cw.write([]string{
			"@THIS",
			"D=M",
		})
	case 1:
		cw.write([]string{
			"@THAT",
			"D=M",
		})
	}
	cw.writePushDRegister()
}

func (cw *CodeWriter) writePushStaticSegment(index int) {
	cw.write([]string{
		fmt.Sprintf("@%s.%d", cw.fileName, index),
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

func (cw *CodeWriter) writePopPointerSegment(index int) {
	cw.write([]string{
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
	})

	switch index {
	case 0:
		cw.write([]string{
			"@THIS",
			"M=D",
		})
	case 1:
		cw.write([]string{
			"@THAT",
			"M=D",
		})
	}
}

func (cw *CodeWriter) writePopStaticSegment(index int) {
	cw.write([]string{
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
	})

	cw.write([]string{
		fmt.Sprintf("@%s.%d", cw.fileName, index),
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

// 図8-4, 8-5 参照
func (cw *CodeWriter) WriteCall(functionName string, numArgs int) {
	returnLabel := fmt.Sprintf("return%d", cw.getReturnLabelNumber())
	cw.write([]string{
		fmt.Sprintf("@%s", returnLabel),
		"D=A",
	})
	// 関数の呼び出し側のLCLなどを保存
	cw.writePushDRegister()

	cw.write([]string{
		"@LCL",
		"D=M",
	})
	cw.writePushDRegister()

	cw.write([]string{
		"@ARG",
		"D=M",
	})
	cw.writePushDRegister()

	cw.write([]string{
		"@THIS",
		"D=M",
	})
	cw.writePushDRegister()

	cw.write([]string{
		"@THAT",
		"D=M",
	})
	cw.writePushDRegister()

	// ARGを別の場所に
	cw.write([]string{
		"@SP",
		"D=M",
		"@5",
		// return label, LCL, ARGなどを格納したため、その分だけSPを戻す
		"D=D-A",
		// さらに関数の引数分だけ戻す。関数の引数は呼び出し時には、既にSPにpushされている。
		fmt.Sprintf("@%d", numArgs),
		"D=D-A",
		// ARG = SP - n - 5 => 既にSPにpushされている引数のはじめのアドレス
		"@ARG",
		"M=D",
		"@SP",
		"D=M",
		// ローカル変数はSPから始まる
		"@LCL",
		"M=D",
		// 関数名ラベルのところにjump
		fmt.Sprintf("@%s", functionName),
		"0;JMP",
		// returnされたらここに戻る
		fmt.Sprintf("(%s)", returnLabel),
	})
}

func (cw *CodeWriter) WriteGoto(label string) {
	cw.write([]string{
		fmt.Sprintf("@%s", cw.getLabelName(label)),
		"0;JMP",
	})
}

func (cw *CodeWriter) WriteIf(label string) {
	cw.write([]string{
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		fmt.Sprintf("@%s", cw.getLabelName(label)),
		"D;JNE",
	})
}

func (cw *CodeWriter) WriteLabel(label string) {
	cw.write([]string{
		fmt.Sprintf("(%s)", cw.getLabelName(label)),
	})
}

func (cw *CodeWriter) WriteFunction(functionName string, numLocals int) {
	cw.write([]string{
		fmt.Sprintf("(%s)", functionName),
		"D=0",
	})

	for i := 0; i < numLocals; i++ {
		cw.writePushDRegister()
	}

	cw.currentFunctionName = functionName
}

func (cw *CodeWriter) WriteReturn() {
	cw.write([]string{
		"@LCL",
		"D=M",
		"@R13", // R13=FRAME
		"M=D",  // FRAME=LCL

		"@5",
		"D=A",
		"@R13",
		"A=M-D", // *(FRAME-5) = *(LCL-5) = return address
		"D=M",
		"@R14",
		"M=D", // R14 = *(FRAME-5) = return address

		"@SP", // pop
		"M=M-1",
		"A=M",
		"D=M", // popした値をDに代入（つまり、戻り値をDに代入）
		"@ARG",
		"A=M",
		"M=D", // 関数実行時はARGだったアドレスに関数の戻り値を代入
		// このあたりは図8-8を参照。

		"@ARG",
		"D=M+1", // ARGの次のアドレス
		"@SP",
		"M=D", // SPの位置をARGの次のアドレスに指定

		"@R13",
		"AM=M-1", // LCL-1 => 呼び出し元のTHATのアドレスが入っている
		"D=M",
		"@THAT",
		"M=D",

		"@R13",
		"AM=M-1", // LCL-2 => 呼び出し元のTHISのアドレスが入っている
		"D=M",
		"@THIS",
		"M=D",

		"@R13",
		"AM=M-1", // LCL-3 => 呼び出し元のARGのアドレスが入っている
		"D=M",
		"@ARG",
		"M=D",

		"@R13",
		"AM=M-1", // LCL-4 => 呼び出し元のLCLのアドレスが入っている
		"D=M",
		"@LCL",
		"M=D",

		"@R14",
		"A=M", // return address
		"0;JMP",
	})
}

func (cw *CodeWriter) getCompLabelNumber() int {
	cw.compNumber++
	return cw.compNumber
}

func (cw *CodeWriter) getReturnLabelNumber() int {
	cw.returnNumber++
	return cw.returnNumber
}

func (cw *CodeWriter) getLabelName(label string) string {
	if cw.currentFunctionName == "" {
		return fmt.Sprintf("MainFunction$%s", label)
	}

	return fmt.Sprintf("%s$%s", cw.currentFunctionName, label)
}

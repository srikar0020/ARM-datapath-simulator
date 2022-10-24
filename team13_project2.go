package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Instruction struct {
	typeofInstruction string
	rawInstruction    string
	linevalue         int64
	programCnt        int64
	opcode            int64
	op                string
	rd                uint8
	rn                uint8
	rm                uint8
	shamt             uint8
	address           uint16
	offset            int64
	immN              int16
	immS              string
}
type Simulation struct {
	r      [32]int32
	cycleN int32
}

var line Instruction
var sim Simulation

func registers() string {
	var registersLine string = "registers:\n"

	for j := 0; j < 4; j++ {
		switch j {
		case 0:
			registersLine += "r00:\t"
		case 1:
			registersLine += "r08:\t"
		case 2:
			registersLine += "r16:\t"
		case 3:
			registersLine += "r32:\t"
		}
		for i := 0; i < 8; i++ {
			registersLine += strconv.FormatInt(int64(sim.r[i+(j*8)]), 10) + "\t"
		}
		registersLine += "\n"
	}
	registersLine += "\n"
	return registersLine
}
func firstLine() string {
	var FLine string = "====================" + "\n\n"
	return FLine
}

func r1Sim() string {
	sim.cycleN++
	var SecondLine string = "cycle:" + strconv.FormatInt(int64(sim.cycleN), 10) + "\t"
	SecondLine += strconv.FormatInt(int64(line.programCnt), 10) + "\t"
	SecondLine += line.op + "\t"
	SecondLine += "	R" + strconv.FormatInt(int64(line.rd), 10)
	SecondLine += ", R" + strconv.FormatInt(int64(line.rn), 10)
	SecondLine += ", R" + strconv.FormatInt(int64(line.rm), 10) + "\n"
	return firstLine() + SecondLine + registers()

}
func r2Sim() string {
	sim.cycleN++
	var SecondLine string = "cycle:" + strconv.FormatInt(int64(sim.cycleN), 10) + "\t"
	SecondLine += strconv.FormatInt(int64(line.programCnt), 10) + "\t"
	SecondLine += line.op + "\t"
	SecondLine += "	R" + strconv.FormatInt(int64(line.rd), 10)
	SecondLine += ", R" + strconv.FormatInt(int64(line.rn), 10)
	SecondLine += ", #" + strconv.FormatInt(int64(line.rm), 10) + "\n"
	return firstLine() + SecondLine + registers()

}

func immediate() string {
	sim.cycleN++
	var SecondLine string = "cycle:" + strconv.FormatInt(int64(sim.cycleN), 10) + "\t"
	SecondLine += strconv.FormatInt(int64(line.programCnt), 10) + "\t"
	SecondLine += line.op + "\t"
	SecondLine += " 	R" + strconv.FormatInt(int64(line.rd), 10)
	SecondLine += ", R" + strconv.FormatInt(int64(line.rn), 10)
	SecondLine += ", #" + strconv.FormatInt(int64(line.immN), 10) + "\n"

	return firstLine() + SecondLine + registers()
}

func main() {
	var InputFileName *string
	var OutputFileName *string
	//var OutputFileNameSim *string
	//var line Instruction
	InputFileName = flag.String("i", "", "Gets the input file name")
	OutputFileName = flag.String("o", "", "Gets the input file name")
	//OutputFileNameSim = flag.String("o", "", "Gets the input file name")
	flag.Parse()
	//fmt.Println("Input:", *InputFileName)
	//fmt.Println("Output:", *OutputFileName)
	fmt.Println(flag.NArg())
	if flag.NArg() != 0 {
		os.Exit(200)
	}
	f, err := os.Create(*OutputFileName + "_dis.txt")
	if err != nil {
		log.Fatal(err)
	}
	s, err := os.Create(*OutputFileName + "_sim.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	file, err := os.Open(*InputFileName)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	file.Close()
	var breakFlag bool = false
	line.programCnt = 96
	for _, eachline := range txtlines {
		//Convert string to Integer

		line.linevalue, err = strconv.ParseInt(eachline, 2, 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		//	Finding opcode by masking and shifting
		line.opcode = (line.linevalue & 0xFFE00000) >> 21
		sim.r[1] = 3
		sim.r[2] = 9

		if !breakFlag {
			// code for R instruction sets
			if line.opcode == 1104 || line.opcode == 1112 || line.opcode == 1360 || line.opcode == 1624 || line.opcode == 1872 {

				line.rd = uint8(line.linevalue & 0x1F)
				line.rn = uint8((line.linevalue & 0x3E0) >> 5)
				line.rm = uint8((line.linevalue & 0x1F0000) >> 16)
				//var line.op string
				if line.opcode == 1104 { // print Instruction
					line.op = "AND"
					sim.r[line.rd] = sim.r[line.rm] & sim.r[line.rn]
				} else if line.opcode == 1112 {
					line.op = "ADD"
					sim.r[line.rd] = sim.r[line.rm] + sim.r[line.rn]
				} else if line.opcode == 1360 {
					line.op = "ORR"
					sim.r[line.rd] = sim.r[line.rm] | sim.r[line.rn]

				} else if line.opcode == 1624 {
					line.op = "SUB"
					sim.r[line.rd] = sim.r[line.rn] - sim.r[line.rm]
				} else if line.opcode == 1872 {
					line.op = "EOR"
					sim.r[line.rd] = sim.r[line.rm] ^ sim.r[line.rn]
				}
				// AND Rd, Rn, Rm

				strArr := []string{eachline[0:11], eachline[11:16], eachline[16:22], eachline[22:27], eachline[27:32]}
				strOut := strings.Join(strArr, " ")
				strOut = strOut + "	" + strconv.FormatInt(int64(line.programCnt), 10)
				strOut = strOut + "	" + line.op
				strOut = strOut + "	R" + strconv.FormatInt(int64(line.rd), 10)
				strOut = strOut + ", R" + strconv.FormatInt(int64(line.rn), 10)
				strOut = strOut + ", R" + strconv.FormatInt(int64(line.rm), 10) + "\n"
				_, err := f.WriteString(strOut)

				if err != nil {
					log.Fatal(err)
				}
				_, err = s.WriteString(r1Sim())

				if err != nil {
					log.Fatal(err)
				}
				line.programCnt = line.programCnt + 4

			}

			// code for R instruction set - LSR, LSL, ASR
			if line.opcode == 1690 || line.opcode == 1691 || line.opcode == 1692 {
				//var line.op string
				line.rd = uint8(line.linevalue & 0x1F)
				line.rn = uint8((line.linevalue & 0x3E0) >> 5)
				line.shamt = uint8((line.linevalue & 0xFC00) >> 10)
				if line.opcode == 1690 {
					line.op = "LSR"
					sim.r[line.rd] = sim.r[line.rn] >> line.shamt
				} else if line.opcode == 1691 {
					line.op = "LSL"
					sim.r[line.rd] = sim.r[line.rn] << line.shamt
				} else if line.opcode == 1692 {
					line.op = "ASR"
					sim.r[line.rd] = sim.r[line.rn] >> line.shamt
				}
				// AND Rd, Rn, Shamt
				strArr := []string{eachline[0:11], eachline[11:16], eachline[16:22], eachline[22:27], eachline[27:32]}
				strOut := strings.Join(strArr, " ")
				strOut = strOut + "	" + strconv.FormatInt(int64(line.programCnt), 10)
				strOut = strOut + "	" + line.op
				strOut = strOut + "	R" + strconv.FormatInt(int64(line.rd), 10)
				strOut = strOut + ", R" + strconv.FormatInt(int64(line.rn), 10)
				strOut = strOut + ", #" + strconv.FormatInt(int64(line.shamt), 10) + "\n"

				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
				_, err = s.WriteString(r2Sim())
				if err != nil {
					log.Fatal(err)
				}
				line.programCnt = line.programCnt + 4
			}

			// Code for D Instruction set
			if line.opcode == 1984 || line.opcode == 1986 {
				//var line.op string
				if line.opcode == 1984 {
					line.op = "STUR"
				} else if line.opcode == 1986 {
					line.op = "LDUR"
				}
				//LDUR Rt, [Rn, #line.programCnt]
				strArr := []string{eachline[0:11], eachline[11:20], eachline[20:22], eachline[22:27], eachline[27:32]}
				strOut := strings.Join(strArr, " ")
				strOut = strOut + "	" + strconv.FormatInt(int64(line.programCnt), 10)
				strOut = strOut + "	" + line.op
				strOut = strOut + "	R" + strconv.FormatInt(int64(line.linevalue&0x1F), 10)
				strOut = strOut + ", [R" + strconv.FormatInt(int64((line.linevalue&0x3E0)>>5), 10)
				strOut = strOut + ", #" + strconv.FormatInt(int64((line.linevalue&0x1FF000)>>12), 10) + "]\n"

				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
				line.programCnt = line.programCnt + 4

			}
			// I format instructions
			if line.opcode == 1160 || line.opcode == 1161 || line.opcode == 1672 || line.opcode == 1673 {
				// ADDI
				// code for 2's compliment
				var x uint16 = uint16((line.linevalue & 0x3FFC00) >> 10)
				var xs int16 = int16((line.linevalue & 0x3FFC00) >> 10)
				if x > 0x7FF {
					x = ^x + 1
					x = x << 4
					x = x >> 4
					xs = int16(x) * -1
				}
				line.immN = xs
				line.rd = uint8(line.linevalue & 0x1F)
				line.rn = uint8((line.linevalue & 0x3E0) >> 5)
				//var line.op string
				if line.opcode == 1160 || line.opcode == 1161 {
					line.op = "ADDI"
					sim.r[line.rd] = sim.r[line.rn] + int32(line.immN)
				} else if line.opcode == 1672 || line.opcode == 1673 {
					line.op = "SUBI"
					sim.r[line.rd] = sim.r[line.rn] - int32(line.immN)
				}

				// ADDI/SUBI Rd, Rn, #immediate

				strArr := []string{eachline[0:10], eachline[10:22], eachline[22:27], eachline[27:32]}

				strOut := strings.Join(strArr, " ")
				strOut = strOut + "  	" + strconv.FormatInt(int64(line.programCnt), 10) + " 	" + line.op
				strOut = strOut + " 	R" + strconv.FormatInt(int64(line.rd), 10)
				strOut = strOut + ", R" + strconv.FormatInt(int64(line.rn), 10)
				strOut = strOut + ", #" + strconv.FormatInt(int64(line.immN), 10) + "\n"

				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
				_, err = s.WriteString(immediate())

				if err != nil {
					log.Fatal(err)
				}
				line.programCnt = line.programCnt + 4
			}

			// for B instruction
			if line.opcode >= 160 && line.opcode <= 191 {
				var bx uint32 = uint32(line.linevalue & 0x3FFFFFF)
				var bxs int64 = int64(line.linevalue & 0x3FFFFFF)
				if bx > 0x1FFFFFF {
					bx = ^bx + 1
					bx = bx << 6
					bx = bx >> 6
					bxs = int64(bx) * -1
				}
				// B #offset
				line.op = "B"
				strArr := []string{eachline[0:6], eachline[6:32]}
				strOut := strings.Join(strArr, " ")
				strOut = strOut + "   	" + strconv.FormatInt(int64(line.programCnt), 10) + " 	" + line.op
				strOut = strOut + " 	#" + strconv.FormatInt(int64(bxs), 10) + "\n"

				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
				line.programCnt = line.programCnt + 4
			}

			// for CB instructions
			if line.opcode >= 1440 && line.opcode <= 1455 {
				//var line.op string
				if line.opcode >= 1440 && line.opcode <= 1447 {
					line.op = "CBZ"
				} else if line.opcode >= 1448 && line.opcode <= 1455 {
					line.op = "CBNZ"
				}
				var cbx uint32 = uint32((line.linevalue & 0xFFFFE0) >> 5)
				var cbxs int32 = int32((line.linevalue & 0xFFFFE0) >> 5)
				if cbx > 0x3FFFFF {
					cbx = ^cbx + 1
					cbx = cbx << 8
					cbx = cbx >> 8
					cbxs = int32(cbx) * -1
				}

				// B #offset
				strArr := []string{eachline[0:8], eachline[8:27], eachline[27:32]}
				strOut := strings.Join(strArr, " ")
				strOut = strOut + "   	" + strconv.FormatInt(int64(line.programCnt), 10) + " 	" + line.op
				strOut = strOut + " 	R" + strconv.FormatInt(line.linevalue&0x1F, 10)
				strOut = strOut + ", #" + strconv.FormatInt(int64(cbxs), 10) + "\n"

				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
				line.programCnt = line.programCnt + 4

			}

			// for IM instructions
			if (line.opcode >= 1684 && line.opcode <= 1687) || (line.opcode >= 1940 && line.opcode <= 1943) {
				//var line.op string
				if line.opcode >= 1684 && line.opcode <= 1687 {
					line.op = "MOVZ"
				} else if line.opcode >= 1940 || line.opcode <= 1943 {
					line.op = "MOVK"
				}
				// B #offset
				strArr := []string{eachline[0:9], eachline[9:11], eachline[11:27], eachline[27:32]}
				strOut := strings.Join(strArr, " ")
				strOut = strOut + "  	" + strconv.FormatInt(int64(line.programCnt), 10) + " 	" + line.op
				strOut = strOut + " 	R" + strconv.FormatInt(line.linevalue&0x1F, 10)
				strOut = strOut + ", " + strconv.FormatInt((line.linevalue&0x1FFFE0)>>5, 10)
				strOut = strOut + ", LSL " + strconv.FormatInt(((line.linevalue&0x600000)>>21)*16, 10) + "\n"

				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
				line.programCnt = line.programCnt + 4
			}

			if line.opcode == 0x0 {
				strOut := "NOP 	" + strconv.FormatInt(int64(line.programCnt), 10)
				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
				line.programCnt = line.programCnt + 4
			}

			if line.opcode == 2038 {
				strOut := eachline + "     " + strconv.FormatInt(int64(line.programCnt), 10) + "	 BREAK\n"
				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
				line.programCnt = line.programCnt + 4
				breakFlag = true
			}
		} else {
			var x uint32 = uint32(line.linevalue)
			var xs int32 = int32(line.linevalue)
			if x > 0x7FFFFFFF {
				x = ^x + 1
				xs = int32(x) * -1
			}
			strOut := eachline + "     " + strconv.FormatInt(int64(line.programCnt), 10) + "	 " + strconv.FormatInt(int64(xs), 10) + "\n"
			_, err := f.WriteString(strOut)
			if err != nil {
				log.Fatal(err)
			}
			line.programCnt = line.programCnt + 4
		}
	}
}
